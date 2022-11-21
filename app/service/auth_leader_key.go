package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_leader_bind"
	"aiyun_cloud_srv/app/model/auth_leader_key"
	"database/sql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

var LeaderkeyService = new(leaderkeyService)

type leaderkeyService struct{}

func (s *leaderkeyService) FindOne(where interface{}) (res *auth_leader_key.Entity, err error) {
	err = auth_leader_key.M.Where(where).Limit(1).Scan(&res)
	return
}

func (s *leaderkeyService) Add(req *auth_leader_key.AddReq, operator_id int) (result sql.Result, err error) {
	result, err = auth_leader_key.M.Data(g.Map{
		"serial_number": req.SerialNumber,
		"code":          req.Code,
		"hospital_id":   req.HospitalId,
		"manager":       req.Manager,
		"public":        req.Public,
		"operator_id":   operator_id,
		"status":        1,
	}).Insert()
	return
}

//修改授权次数
func (s *leaderkeyService) IncAuthCount(where interface{}) (err error) {
	_, err = auth_leader_key.M.Where(where).Data(g.Map{"auth_count": gdb.Raw("auth_count+1")}).Update()
	return
}

//查询绑定情况
func (s *leaderkeyService) GetBindInfo(where interface{}) (res *auth_leader_bind.Entity, err error) {
	err = auth_leader_bind.M.Where(where).Limit(1).Scan(&res)
	return
}

//添加绑定
func (s *leaderkeyService) Bind(req *auth_leader_key.BindLeaderReq, leader_key *auth_leader_key.Entity, serial_number string) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//记录授权次数
	//_, err = tx.Model(auth_leader_key.Table).Where(g.Map{"id": leader_key.Id}).Data(g.Map{
	//	"auth_count": gdb.Raw("auth_count+1"),
	//}).Update()
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	//添加授权记录
	_, err = tx.Model(auth_leader_bind.Table).Data(g.Map{
		"auth_leader_id":       leader_key.Id,
		"serial_number":        serial_number,
		"hospital_id":          leader_key.HospitalId,
		"client_author_number": req.AuthorNumber,
		"server_author_number": req.ServerAuthorNumber,
		"status":               1,
	}).Insert()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//添加绑定
func (s *leaderkeyService) Rebind(req *auth_leader_key.BindLeaderReq, leader_key *auth_leader_key.Entity, serial_number string) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//记录授权次数
	//_, err = tx.Model(auth_leader_key.Table).Where(g.Map{"id": leader_key.Id}).Data(g.Map{
	//	"auth_count": gdb.Raw("auth_count+1"),
	//}).Update()
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	// 删除之前的绑定记录
	_, err = tx.Model(auth_leader_bind.Table).
		Where("auth_leader_id", leader_key.Id).
		Data(g.Map{"status": 0}).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	//添加授权记录
	_, err = tx.Model(auth_leader_bind.Table).Data(g.Map{
		"auth_leader_id":       leader_key.Id,
		"serial_number":        serial_number,
		"hospital_id":          leader_key.HospitalId,
		"client_author_number": req.AuthorNumber,
		"server_author_number": req.ServerAuthorNumber,
		"status":               1,
	}).Insert()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//添加绑定
func (s *leaderkeyService) Unbind(leader_id int) (err error) {
	_, err = auth_leader_bind.M.
		Where("auth_leader_id", leader_id).
		Data(g.Map{"status": 0}).Update()
	return
}

//设置数据状态
func (s *leaderkeyService) SetStatus(req *model.SetStatusParams) (err error) {
	_, err = auth_leader_key.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *leaderkeyService) Delete(req *model.Ids) (err error) {
	_, err = auth_leader_key.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取分页列表
func (s *leaderkeyService) Page(req *model.PageReqParams) (total int, list []*auth_leader_key.Entity, err error) {
	M := auth_leader_key.M_alias

	if req.KeyWord != "" {
		M = M.WhereOrLike("alk.code", "%"+req.KeyWord+"%")
		M = M.WhereOrLike("alk.serial_number", "%"+req.KeyWord+"%")
		M = M.WhereOrLike("h.name", "%"+req.KeyWord+"%")
	}

	if req.Status != -1 {
		M = M.Where("alk.status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("alk.create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("alk.create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order("alk." + req.Order + " " + req.Sort)
	}

	total, err = M.
		LeftJoin("hospital h", "h.id = alk.hospital_id").
		LeftJoin("auth_leader_bind alb", "alb.auth_leader_id = alk.id AND alb.status = 1").
		Group("alk.id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("alk.*,h.name as hospital_name,alb.server_author_number").
		LeftJoin("hospital h", "h.id = alk.hospital_id").
		LeftJoin("auth_leader_bind alb", "alb.auth_leader_id = alk.id AND alb.status = 1").
		All()

	g.Dump(data)

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*auth_leader_key.Entity, len(data))
	err = data.Structs(&list)
	return
}
