package Admin

import (
	_const "aiyun_cloud_srv/app/const"
	"aiyun_cloud_srv/app/model/sys_admin"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/cache"
	"aiyun_cloud_srv/library/utils/jwt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"time"
)

var Public = publicApi{}

type publicApi struct{}

func (*publicApi) Test(r *ghttp.Request) {
	service.RegionService.Tree()
}

// Index is a demonstration route handler for output "Hello World!".
func (*publicApi) Login(r *ghttp.Request) {
	var req *sys_admin.LoginReq
	//获取参数&参数校验
	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//校验验证码
	//cache_code, err := cache.Get(_const.LOGIN_CODE_CACHE_KEY(req.Username, req.Code))
	//if err != nil {
	//	response.ErrorSys(r, err)
	//}
	//if req.Code != cache_code {
	//	response.Json(r, 2, "验证码错误")
	//}

	adminInfo, err := service.SysAdminService.GetByUsername(req.Username)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if adminInfo == nil {
		response.Error(r, "用户名不存在")
	}

	if adminInfo.Password != utils.GenPwd(req.Password, adminInfo.Salt) {
		response.Error(r, "密码错误")
	}

	//生成token
	token, err := jwt.GenerateLoginToken(adminInfo)
	if err != nil {
		response.ErrorSys(r, err)
	}

	//保存token
	token_expires := g.Cfg().GetInt("jwt.expires")
	if err := cache.Set(_const.TOKEN_CACHE_KEY(token), token, time.Second*time.Duration(token_expires)); err != nil {
		response.ErrorSys(r, err)
	}

	adminInfo.Token = token
	adminInfo.Expires = token_expires

	role_info, err := service.SysRoleService.FindById(adminInfo.RoleId)
	if err != nil {
		response.ErrorDb(r, err)
	}
	adminInfo.RoleRule = role_info.Rule

	response.Success(r, adminInfo)
}

//获取验证码
func (*publicApi) Captcha(r *ghttp.Request) {
	response.Success(r)
}

func (*publicApi) GetRegion(r *ghttp.Request) {
	tree, err := service.RegionService.Tree()

	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}
