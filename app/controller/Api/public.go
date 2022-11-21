package Api

import (
	"aiyun_cloud_srv/app/model/hospital"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/net/ghttp"
)

var Public = publicApi{}

type publicApi struct{}

//获取地区列表数据 is_tree为ture 返回树结构数据
func (*publicApi) Region(r *ghttp.Request) {
	is_tree := r.GetQueryBool("is_tree")
	if is_tree {
		list, err := service.RegionService.Tree()
		if err != nil {
			response.ErrorSys(r, err)
		}
		response.Success(r, list)
	} else {
		list, err := service.RegionService.List()
		if err != nil {
			response.ErrorSys(r, err)
		}
		response.Success(r, list)
	}
}

//通过地区检索医院数据
func (*publicApi) Hospital(r *ghttp.Request) {
	region_id := r.GetQueryString("region_id", 0)

	_, list, err := service.HospitalService.Page(&hospital.PageReqParams{
		RegionId: []string{"", "", region_id},
		PageSize: 1000,
		Status:   1,
	})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}
