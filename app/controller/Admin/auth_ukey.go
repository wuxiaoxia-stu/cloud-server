package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_ukey"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var AuthUKey = authUKeyApi{}

type authUKeyApi struct{}

func (*authUKeyApi) Info(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")

	info, err := service.AuthUkeyService.FindOne(g.Map{"serial_number": serial_number})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, info)
}

// 获取ukey列表
func (*authUKeyApi) List(r *ghttp.Request) {
	var req *auth_ukey.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.AuthUkeyService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加ukey
func (*authUKeyApi) Add(r *ghttp.Request) {
	var req *auth_ukey.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.AuthUkeyService.Add(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//设置ukey状态
func (*authUKeyApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.AuthUkeyService.SetStatus(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除ukey
func (*authUKeyApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.AuthUkeyService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}
