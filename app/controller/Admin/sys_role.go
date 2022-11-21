package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/sys_role"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var SysRole = sysRoleApi{}

type sysRoleApi struct{}

// 获取系统用户列表
func (*sysRoleApi) List(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.SysRoleService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加管理员
func (*sysRoleApi) Add(r *ghttp.Request) {
	var req *sys_role.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.SysRoleService.Add(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//修改管理员信息
func (*sysRoleApi) Edit(r *ghttp.Request) {
	var req *sys_role.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.SysRoleService.Edit(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "编辑成功")
}

//设置管理员状态
func (*sysRoleApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.SysRoleService.SetStatus(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除管理员
func (*sysRoleApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.SysRoleService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//获取权限菜单
func (*sysRoleApi) AuthMenu(r *ghttp.Request) {
	response.Success(r, service.AuthMenuService.Tree())
}
