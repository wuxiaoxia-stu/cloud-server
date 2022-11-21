package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/sys_role"
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

var SysRoleService = new(sysRoleService)

type sysRoleService struct{}

// 查询单个信息
func (s *sysRoleService) FindById(id int) (res *sys_role.Entity, err error) {
	err = sys_role.M.Where("id=?", id).Scan(&res)
	return
}

//添加新角色
func (s *sysRoleService) Add(req *sys_role.AddReq) (err error) {
	_, err = sys_role.M.Data(g.Map{
		"name":   req.Name,
		"rule":   req.Rule,
		"status": 1,
	}).Insert()
	return
}

//更新用户信息
func (s *sysRoleService) Edit(req *sys_role.EditReq) (err error) {
	_, err = sys_role.M.Where("id", req.Id).Data(g.Map{
		"rule": req.Rule,
	}).Update()
	return
}

//设置数据状态
func (s *sysRoleService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = sys_role.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *sysRoleService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = sys_role.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取分页列表
func (s *sysRoleService) Page(req *model.PageReqParams) (total int, list []*sys_role.Entity, err error) {
	M := sys_role.M

	if req.KeyWord != "" {
		M = M.WhereLike("name", "%"+req.KeyWord+"%")
	}

	if req.Status != -1 {
		M = M.Where("status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order(req.Order + " " + req.Sort)
	}

	total, err = M.Group("id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*sys_role.Entity, len(data))
	err = data.Structs(&list)
	return
}
