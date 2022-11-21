package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"aiyun_cloud_srv/app/model/kl_version"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/encrypt"
	"bufio"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/v2/errors/gerror"
	"os"
	"path/filepath"
	"strings"
)

var KlVersionService = new(klVersionService)

type klVersionService struct{}

//获取列表数据
func (s *klVersionService) Info(where interface{}) (res *kl_version.Entity, err error) {
	err = kl_version.M.Where(where).Scan(&res)
	return
}

//获取列表数据
func (s *klVersionService) Page(req *model.PageReqParams) (total int, list []*kl_version.Entity, err error) {
	M := kl_version.M_alias

	if req.KeyWord != "" {
		M = M.WhereLike("kv.name", "%"+req.KeyWord+"%")
	}

	if req.Status != -1 {
		M = M.Where("kv.status=?", req.Status)
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("kv." + req.Order + " " + req.Sort)
	}

	total, err = M.Group("kv.id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("kv.*,sa.username as operator_name").LeftJoin("sys_admin sa", "sa.id = kv.operator_id").All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*kl_version.Entity, len(data))
	err = data.Structs(&list)

	base_url := g.Cfg().GetString("server.Domain")
	for _, v := range list {
		v.DataFullPath = base_url + strings.TrimLeft(v.DataPath, "public/")
	}
	return
}

func (s *klVersionService) Add(req *kl_version.AddReq, operator_id int) (err error) {
	//数据库数据导出为json
	if err = s.ExportFeatureData(); err != nil {
		return err
	}

	if err = s.ExportFeatureAtlasData(); err != nil {
		return err
	}

	if err = s.ExportSyndromeData(); err != nil {
		return err
	}

	if err = s.ExportSyndromeFeatureData(); err != nil {
		return err
	}

	file_list_all, err := utils.GetAllFile("public/attachment/feature_ledge")
	if err != nil {
		return err
	}

	//去除后缀为data的文件
	file_list := []string{}
	for _, v := range file_list_all {
		if !(strings.Contains(v, ".data") || strings.Contains(v, ".zip")) {
			file_list = append(file_list, v)
		}
	}

	//对目录下所有文件加密，生成后缀为.data的加密文件
	encrypt_file_list := []string{}
	for _, v := range file_list {
		dir, fname := filepath.Split(v)
		des := filepath.Join(dir, fname+".data")
		//文件加密
		err = encrypt.EncryptFile(v, des, "123456")
		if err != nil {
			return err
		}
		encrypt_file_list = append(encrypt_file_list, des)
	}

	g.Dump(encrypt_file_list)
	// 保留原来文件的结构 把加密为文件加入压缩包
	err = utils.ZipFiles("public/kl/version/"+req.Name+".zip", encrypt_file_list,
		"public\\attachment\\feature_ledge", "feature_ledge")
	if err != nil {
		return err
	}
	_, err = kl_version.M.Data(g.Map{
		"name":        req.Name,
		"operator_id": operator_id,
		"data_path":   "public/kl/version/" + req.Name + ".zip",
	}).Insert()
	return
}

//设置数据状态
func (s *klVersionService) SetStatus(req *model.SetStatusParams) (err error) {
	_, err = kl_version.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *klVersionService) Delete(req *model.Ids) (err error) {
	_, err = kl_version.M.WhereIn("id", req.Ids).Delete()
	return
}

func (s *klVersionService) ExportFeatureData() (err error) {
	kl_feature_list := []*kl_feature.Entity{}
	err = kl_feature.M.Scan(&kl_feature_list)
	if err != nil {
		return err
	}
	b, err := json.Marshal(kl_feature_list)
	if err != nil {
		return nil
	}
	file, err := os.OpenFile("public/attachment/feature_ledge/data/kl_feature.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return writer.Flush()
}

func (s *klVersionService) ExportFeatureAtlasData() (err error) {
	kl_feature_atlas_list := []*kl_feature_atlas.Entity{}
	err = kl_feature_atlas.M.Scan(&kl_feature_atlas_list)
	if err != nil {
		return err
	}
	b, err := json.Marshal(kl_feature_atlas_list)
	if err != nil {
		return nil
	}
	file, err := os.OpenFile("public/attachment/feature_ledge/data/kl_feature_atlas.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return writer.Flush()
}

func (s *klVersionService) ExportSyndromeData() (err error) {
	kl_syndrome_list := []*kl_syndrome.Entity{}
	err = kl_syndrome.M.Scan(&kl_syndrome_list)
	if err != nil {
		return err
	}
	b, err := json.Marshal(kl_syndrome_list)
	if err != nil {
		return nil
	}
	file, err := os.OpenFile("public/attachment/feature_ledge/data/kl_syndrome.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return writer.Flush()
}

func (s *klVersionService) ExportSyndromeFeatureData() (err error) {
	kl_syndrome_feature_list := []*kl_syndrome_feature.Entity{}
	err = kl_syndrome_feature.M.Scan(&kl_syndrome_feature_list)
	if err != nil {
		return err
	}
	b, err := json.Marshal(kl_syndrome_feature_list)
	if err != nil {
		return nil
	}
	file, err := os.OpenFile("public/attachment/feature_ledge/data/kl_syndrome_feature.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(b); err != nil {
		return err
	}

	return writer.Flush()
}
