package Api

import (
	"aiyun_cloud_srv/app/model/hospital"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
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

	_, list, err := service.HospitalService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}
