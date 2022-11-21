package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/sys_admin"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var SysAdmin = sysAdminApi{}

type sysAdminApi struct{}

// 获取系统用户列表
func (*sysAdminApi) List(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.SysAdminService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加管理员
func (*sysAdminApi) Add(r *ghttp.Request) {
	var req *sys_admin.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	bool, err := service.SysAdminService.CheckDepartmentId(req.DepartmentId)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if !bool {
		response.Error(r, "参数错误，部门信息异常")
	}

	_, err = service.SysAdminService.Add(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//修改管理员信息
func (*sysAdminApi) Edit(r *ghttp.Request) {
	var req *sys_admin.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	bool, err := service.SysAdminService.CheckDepartmentId(req.DepartmentId)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if !bool {
		response.Error(r, "参数错误，部门信息异常")
	}

	_, err = service.SysAdminService.Edit(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "编辑成功")
}

//设置管理员状态
func (*sysAdminApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.Id == 1 {
		response.Error(r, "系统级管理员不能被禁止")
	}

	_, err := service.SysAdminService.SetStatus(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除管理员
func (*sysAdminApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if utils.InArray(1, req.Ids) {
		response.Error(r, "系统级管理员禁止删除")
	}

	_, err := service.SysAdminService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//获取部门结构数据
func (*sysAdminApi) GetDepartmentTree(r *ghttp.Request) {
	tree, err := service.SysAdminService.GetDepartmentTree()
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}

func (*sysAdminApi) GetRoleRule(r *ghttp.Request) {
	role_id := r.GetCtxVar("role_id").Int()
	role, err := service.SysRoleService.FindById(role_id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, role.Rule)
}

//获取路由列表
func (*sysAdminApi) GetAsyncRoutes(r *ghttp.Request) {

	response.Success(r, g.Array{})

	role_id := r.GetCtxVar("role_id").Int()

	role, err := service.SysRoleService.FindById(role_id)

	if err != nil {
		response.ErrorDb(r, err)
	}
	if role.Status != 1 {
		response.Success(r, g.Array{})
	}

	if role.Rule == "[]" {
		response.Success(r, g.Array{})
	}

	var rule []string
	if err := json.Unmarshal([]byte(role.Rule), &rule); err != nil {
		response.ErrorSys(r, err)
	}

	menu_all := service.AuthMenuService.Tree()

	menu := g.Array{}
	//for i, v := range menu_all {
	//	exist := false
	//	if utils.StrInArray(v.Api, rule){
	//		exist = true
	//	}
	//
	//	if exist {
	//		child_arr := g.Array{}
	//		for i2, v2 := range v.Children {
	//			exist := false
	//			if utils.StrInArray(v.Api, rule){
	//				exist = true
	//			}
	//
	//			if exist {
	//				child_arr = append(child_arr, g.Map{
	//					"path": v2.Path,
	//					"meta": g.Map{
	//						"title": v2.Label,
	//						"icon":  v2.Icon,
	//						"rank":  i2 + 1,
	//					},
	//					"children": g.Array{},
	//				})
	//			}
	//		}
	//
	//		menu = append(menu, g.Map{
	//			"path": v.Path,
	//			"redirect": v.Redirect,
	//			"meta": g.Map{
	//				"title": v.Label,
	//				"icon":  v.Icon,
	//				"rank":  i + 1,
	//			},
	//			"children": child_arr,
	//		})
	//	}
	//
	//}
	//
	//response.Success(r, menu)
	//
	for i, v := range menu_all {
		children := g.Array{}

		for i2, v2 := range v.Children {
			child_exist := false
			for _, v3 := range v2.Children {
				if utils.StrInArray(v3.Api, rule) {
					child_exist = true
				}
			}

			if child_exist {
				children = append(children, g.Map{
					"path": v2.Path,
					"meta": g.Map{
						"title": v2.Label,
						"icon":  v2.Icon,
						"rank":  i2,
					},
				})
			}
		}

		if len(children) > 0 {
			menu = append(menu, g.Map{
				"path":     v.Path,
				"redirect": v.Redirect,
				"meta": g.Map{
					"title": v.Label,
					"icon":  v.Icon,
					"rank":  i,
				},
				"children": children,
			})
		}
	}
	response.Success(r, menu)
}
