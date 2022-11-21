package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/hospital"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Hospital = hospitalApi{}

type hospitalApi struct{}

// 获取医院列表
func (*hospitalApi) List(r *ghttp.Request) {
	var req *hospital.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.HospitalService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加医院
func (*hospitalApi) Add(r *ghttp.Request) {
	var req *hospital.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.HospitalService.Add(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

func (*hospitalApi) Edit(r *ghttp.Request) {
	var req *hospital.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.HospitalService.Edit(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//设置医院状态
func (*hospitalApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.HospitalService.SetStatus(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除医院
func (*hospitalApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.HospitalService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}
