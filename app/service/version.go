package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/version"
	"github.com/gogf/gf/frame/g"
)

var VersionService = new(versionService)

type versionService struct{}

func (s *versionService) List(where interface{}) (list []*version.Entity, err error) {
	base_url := g.Cfg().GetString("server.Domain")
	data, err := version.M_alias.
		Fields("v.*,sa.username as operator_name,CONCAT('"+base_url+"',package_url) as package_full_path").
		LeftJoin("sys_admin sa", "sa.id = v.operator_id").
		Where(where).
		Order("v.id DESC").
		All()
	if err != nil {
		return
	}
	list = make([]*version.Entity, len(data))
	err = data.Structs(&list)
	return
}

//添加新版本
func (s *versionService) Add(req *version.Entity) (err error) {
	_, err = version.M.Data(g.Map{
		"channel":        req.Channel,
		"version_number": req.VersionNumber,
		"info":           req.Info,
		"bug_fix":        req.BugFix,
		"update_range":   req.UpdateRange,
		"package_url":    req.PackageUrl,
		"remark":         req.Remark,
		"operator_id":    req.OperatorId,
		"status":         1,
	}).Insert()
	return
}

//更新用户信息
func (s *versionService) Edit(req *version.Entity) (err error) {
	_, err = version.M.Data(g.Map{
		"channel":        req.Channel,
		"version_number": req.VersionNumber,
		"info":           req.Info,
		"bug_fix":        req.BugFix,
		"update_range":   req.UpdateRange,
		"package_url":    req.PackageUrl,
		"remark":         req.Remark,
		"operator_id":    req.OperatorId,
		"status":         req.Status,
	}).Where("id", req.Id).Update()
	return
}

//设置数据状态
func (s *versionService) SetStatus(req *model.SetStatusParams) (err error) {
	_, err = version.M.WhereIn("id", req.Id).Data(g.Map{"status": req.Status}).Update()
	return
}

//批量删除
func (s *versionService) Delete(req *model.Ids) (err error) {
	_, err = version.M.WhereIn("id", req.Ids).Delete()
	return
}
