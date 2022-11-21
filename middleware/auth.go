package middleware

import (
	_const "aiyun_cloud_srv/app/const"
	"aiyun_cloud_srv/app/model/sys_admin"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/cache"
	"aiyun_cloud_srv/library/utils/jwt"
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

// Auth 权限判断处理中间件
func Auth(r *ghttp.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		response.Json(r, 9, "Token错误")
	}
	token, err := cache.Get(_const.TOKEN_CACHE_KEY(tokenStr))
	if err != nil || token == "" {
		response.Json(r, 9, "Token无效或已过期", g.Array{})
	}

	data, err := jwt.ParseToken(tokenStr, []byte(g.Cfg().GetString("jwt.sign", "qc_sign")))
	if err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}

	var u *sys_admin.Entity
	if err = gconv.Struct(data, &u); err != nil {
		g.Log().Error(err.Error())
		response.Json(r, 9, "Token错误")
	}
	r.SetCtxVar("uid", u.Id)
	r.SetCtxVar("role_id", u.RoleId)

	//判断角色权限
	is_allow, err := RoleAuth(u.RoleId, strings.TrimPrefix(r.Router.Uri, "/admin/"))
	if err != nil {
		response.ErrorSys(r, err)
	}

	if !is_allow {
		response.Json(r, 7, "权限不足")
	}

	// 执行下一步请求逻辑
	r.Middleware.Next()
}

//角色白名单
var RoleAuthWhiteList = []string{
	"file/upload",
	"index/order-region-data",
	"index/device-data",
	"order/meal-options",
	"kl-feature/group-list",
	"kl-syndrome/type",
	"hospital/list",
	"sys-admin/get-department-tree",
	"sys-admin/get-async-routes",
}

//获取角色菜单数据
func RoleAuth(role_id int, uri string) (is_allow bool, err error) {
	if utils.StrInArray(uri, RoleAuthWhiteList) {
		return true, nil
	}

	role, err := service.SysRoleService.FindById(role_id)
	if err != nil {
		return
	}

	if role == nil {
		return false, gerror.New("角色不存在")
	}

	rule := []string{}
	err = json.Unmarshal([]byte(role.Rule), &rule)
	if err != nil {
		return
	}

	if utils.StrInArray(uri, rule) {
		is_allow = true
	}

	return
}
