package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_licence"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Licence = licenceApi{}

type licenceApi struct{}

func (*licenceApi) List(r *ghttp.Request) {
	var req *auth_licence.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.AuthLicenceService.Page(req)

	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//设置状态
func (*licenceApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if req.Status == 1 {
		b, err := service.AuthLicenceService.CheckNewLicence(req.Id)
		if err != nil {
			response.ErrorDb(r, err)
		}
		if b {
			response.Error(r, "此设备存在新的授权证书，禁止此操作")
		}
	}

	if err := service.AuthLicenceService.SetStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除
func (*licenceApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.AuthLicenceService.Delete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}
