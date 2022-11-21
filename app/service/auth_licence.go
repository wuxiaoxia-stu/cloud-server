package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_client"
	"aiyun_cloud_srv/app/model/auth_licence"
	"aiyun_cloud_srv/app/model/auth_licence_num"
	"aiyun_cloud_srv/app/model/auth_ukey"
	"aiyun_cloud_srv/app/model/order"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/rsa_crypt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"strings"
	"time"
)

var AuthLicenceService = new(authLicenceService)

type authLicenceService struct{}

//查询数据
func (s *authLicenceService) Find(where g.Map) (list []*auth_licence.Entity, err error) {
	err = auth_licence.M.Where(where).Scan(&list)
	return
}

func (s *authLicenceService) FindOne(where g.Map) (list *auth_licence.Entity, err error) {
	err = auth_licence.M.Where(where).Limit(1).Scan(&list)
	return
}

//查询全部信息
func (s *authLicenceService) InfoAll(device_serial_number string) (res *auth_licence.Entity, err error) {
	err = auth_licence.M.Where("device_serial_number", device_serial_number).Scan(&res)
	if err != nil || res == nil {
		return
	}

	err = auth_client.M.Where("id", res.AuthClientId).Scan(&res.ClientInfo)
	if err != nil {
		return
	}
	return
}

//查询授权是否存在
func (s *authLicenceService) AuthExist(uuid string, role int) (auth_licence_info *auth_licence.Entity, err error) {
	err = auth_licence.M.Where(g.Map{"uuid": uuid, "role": role, "status": 1}).Order("id DESC").Scan(&auth_licence_info)
	return
}

func (s *authLicenceService) Create(req *auth_client.Entity, ukey_code string, o *order.Entity) (res *auth_licence.Entity, err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// 添加客户端信息记录
	if _, err = tx.Model(auth_client.Table).Data(g.Map{
		"disk_id":            req.DiskID,
		"cpu_name":           req.CpuName,
		"cpu_processor_id":   req.CpuID,
		"baseboard_id":       req.BaseboardID,
		"uuid":               req.Uuid,
		"gpu_name":           req.GpuName,
		"role":               req.Role,
		"ukey_serial_number": req.UkeySerialNumber,
		"client_version":     req.ClientVersion,
		"ai_version":         req.AiVersion,
		"status":             1,
	}).Insert(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//生成授权编号
	author_number, err := s.GenAuthorNumber(req.Role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//生成公私钥对
	rsa_pair := rsa_crypt.RsaKeyPair{}
	if rsa_crypt.RSAGenKey(2048, &rsa_pair); err != nil {
		tx.Rollback()
		return nil, err
	}

	var client_info *auth_client.Entity
	if err = tx.Model(auth_client.Table).Where(g.Map{"uuid": req.Uuid, "role": req.Role}).Scan(&client_info); err != nil {
		tx.Rollback()
		return nil, err
	}

	licence, err := s.genLicenceSignature(g.Map{
		"author_number":      author_number,
		"client_version":     req.ClientVersion,
		"role":               req.Role,
		"hospital_id":        req.HospitalId,
		"serial_number":      req.SerialNumber,
		"ukey_serial_number": req.UkeySerialNumber,
		"uuid":               req.Uuid,
	}, rsa_pair.PriKey)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//添加授权记录
	if _, err = auth_licence.M.Data(g.Map{
		"auth_client_id":       client_info.Id,
		"uuid":                 req.Uuid,
		"role":                 req.Role,
		"author_number":        author_number,
		"hospital_id":          req.HospitalId,
		"licence":              licence,
		"public_key":           rsa_pair.PubKey,
		"private_key":          rsa_pair.PriKey,
		"ukey_code":            ukey_code,
		"ukey_serial_number":   req.UkeySerialNumber,
		"device_serial_number": req.SerialNumber,
	}).Insert(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//增加ukey使用次数
	if _, err = tx.Model(auth_ukey.Table).
		Where("code", ukey_code).
		Data(g.Map{"used_times": gdb.Raw("used_times+1")}).
		Update(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//客户端授权
	if req.Role == 1 {
		// 修改授权时间
		if _, err = tx.Model(order_device.Table).
			Where("serial_number", req.SerialNumber).
			Data(g.Map{"authorize_time": time.Now().Unix()}).
			Update(); err != nil {
			tx.Rollback()
			return nil, err
		}

		//修改订单状态
		count, err := tx.Model(order_device.Table).Where(g.Map{
			"order_number":     o.OrderNumber,
			"authorize_time >": 0,
			"status":           1,
		}).Count()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if count >= o.Count {
			_, err := tx.Model(order.Table).Where(g.Map{"order_number": o.OrderNumber}).Data(g.Map{"status": 1}).Update()
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	res = &auth_licence.Entity{
		AuthClientId:     client_info.Id,
		AuthorNumber:     author_number,
		UkeySerialNumber: req.UkeySerialNumber,
		PublicKey:        rsa_pair.PubKey,
		PrivateKey:       rsa_pair.PriKey,
		Licence:          licence,
	}

	return res, tx.Commit()
}

func (s *authLicenceService) genLicenceSignature(signMap map[string]interface{}, priKeyBase64 string) (res string, err error) {
	hashStr := strings.Join(utils.Map2List(signMap), "&")
	hashStr = utils.EncodeSHA256(hashStr, "")
	priKey, err := rsa_crypt.LoadPrivateKeyBase64(priKeyBase64)
	if err != nil {
		return
	}

	signMap["signature"] = rsa_crypt.Sign(priKey, hashStr)

	delete(signMap, "uuid")

	signMapByte, err := json.Marshal(signMap)
	if err != nil {
		return
	}

	res = base64.StdEncoding.EncodeToString(signMapByte)
	return
}

//生成服务端或客户端授权编号
func (s *authLicenceService) GenAuthorNumber(role int) (author_number string, err error) {

	var res *auth_licence_num.Entity

	now := time.Now().Format("0601")

	err = auth_licence_num.M.Where("l_month", now).Scan(&res)
	if err != nil {
		return
	}

	number := 1
	if res != nil { //更新number
		number = res.Number + 1
		_, err = auth_licence_num.M.Where("l_month", now).Data(g.Map{
			"number": number,
		}).Update()
	} else { //添加新的记录
		_, err = auth_licence_num.M.Data(g.Map{
			"l_month": now,
			"number":  number,
		}).Insert()

	}

	if role == 1 { //C端
		author_number = fmt.Sprintf("C_%s%04d", now, number)
	} else { //S端
		author_number = fmt.Sprintf("S_%s%04d", now, number)
	}
	return
}

//获取分页列表
func (s *authLicenceService) Page(req *auth_licence.PageReqParams) (total int, list []*auth_licence.Entity, err error) {
	M := auth_licence.M_alias

	if req.AuthorNumber != "" {
		M = M.WhereLike("al.author_number", "%"+req.AuthorNumber+"%")
	}

	if req.SerialNumber != "" {
		M = M.WhereLike("al.device_serial_number", "%"+req.SerialNumber+"%")
	}

	if req.Uuid != "" {
		M = M.WhereLike("ac.uuid", "%"+req.Uuid+"%")
	}

	if req.UkeyCode != "" {
		M = M.Where("al.ukey_code", req.UkeyCode)
	}

	if req.HospitalId > 0 {
		M = M.Where("h.id", req.HospitalId)
	} else if len(req.RegionIds) > 0 {
		M = M.WhereLike("h.region_id", "%,"+strings.Join(req.RegionIds, ","))
	}

	if req.Status != -1 {
		M = M.Where("al.status=?", req.Status)
	}

	if req.Role > 0 {
		M = M.Where("al.role", req.Role)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("al.create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("al.create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order("al." + req.Order + " " + req.Sort)
	} else {
		M = M.Order("al.id DESC")
	}

	total, err = M.
		LeftJoin("auth_client ac", "ac.id = al.auth_client_id").
		LeftJoin("hospital h", "h.id = al.hospital_id").
		Group("al.id").
		Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("al.id,al.licence_id,al.auth_client_id,al.role,al.author_number,al.ukey_code,al.ukey_serial_number,"+
		"al.device_serial_number,al.update_at,al.pair_at,al.status,"+
		"ac.uuid,ac.client_version,ac.ai_version,h.name as hospital_name,al2.author_number as server_author_number").
		LeftJoin("auth_licence al2", "al.licence_id = al2.id").
		LeftJoin("auth_client ac", "ac.id = al.auth_client_id").
		LeftJoin("hospital h", "h.id = al.hospital_id").
		All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*auth_licence.Entity, len(data))
	err = data.Structs(&list)

	id_arr := []int{}
	for _, v := range list {
		if v.Role == 2 {
			id_arr = append(id_arr, v.Id)
		}
	}

	var list_ []*auth_licence.Entity
	auth_licence.M.WhereIn("licence_id", id_arr).Scan(&list_)
	for _, v := range list {
		for _, v2 := range list_ {
			if v.Id == v2.LicenceId {
				v.PairCount += 1
				v.PairClientAuthorNumber = append(v.PairClientAuthorNumber, v2.AuthorNumber)
			}
		}
	}

	return
}

// 检查是否有新的授权证书
func (s *authLicenceService) CheckNewLicence(id int) (b bool, err error) {
	var info *auth_licence.Entity
	err = auth_licence.M.Where("id", id).Where("status", 1).Scan(&info)
	if err != nil {
		return
	}

	if info != nil {
		var new_info *auth_licence.Entity
		err = auth_licence.M.Where(g.Map{
			"id !=":         info.Id,
			"author_number": info.AuthorNumber,
			"status":        1,
		}).Scan(&new_info)
		if err != nil {
			return
		}

		if new_info != nil {
			b = true
		}
	}
	return
}

//设置数据状态
func (s *authLicenceService) SetStatus(req *model.SetStatusParams) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	var licence_info *auth_licence.Entity
	if err := tx.Model(auth_licence.Table).Where("id", req.Id).Scan(&licence_info); err != nil {
		tx.Rollback()
		return err
	}

	if licence_info == nil {
		return fmt.Errorf("数据不存在")
	}
	if _, err = tx.Model(auth_licence.Table).Where("id", req.Id).Data(req).Update(); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Model(auth_client.Table).Where("id", licence_info.AuthClientId).Data(g.Map{"status": req.Status}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//批量删除
func (s *authLicenceService) Delete(req *model.Ids) (err error) {
	_, err = auth_licence.M.WhereIn("id", req.Ids).Delete()
	return
}
