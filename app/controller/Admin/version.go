package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/version"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Version = versionApi{}

type versionApi struct{}

//列表
func (*versionApi) List(r *ghttp.Request) {
	keywords := r.GetQueryString("keywords")
	channel := r.GetQueryInt("channel")
	status := r.GetQueryInt("status")

	where := g.Map{}
	if keywords != "" {
		where["v.version_number"] = "%" + keywords + "%"
	}

	if channel > 0 {
		where["v.channel"] = channel
	}

	if status > 0 {
		where["v.status"] = status
	}

	list, err := service.VersionService.List(where)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

//添加版本
func (*versionApi) Add(r *ghttp.Request) {
	var req *version.Entity

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.VersionService.Add(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//修改版本
func (*versionApi) Edit(r *ghttp.Request) {
	var req *version.Entity

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.VersionService.Edit(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

//删除版本
func (*versionApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.VersionService.Delete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//设置状态
func (*versionApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.VersionService.SetStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "状态设置成功")
}
