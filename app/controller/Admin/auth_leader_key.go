package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_leader_key"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var LeaderKey = leaderKeyApi{}

type leaderKeyApi struct{}

func (*leaderKeyApi) Info(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")

	info, err := service.LeaderkeyService.FindOne(g.Map{"serial_number": serial_number})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, info)
}

// 获取主任秘钥列表
func (*leaderKeyApi) List(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.LeaderkeyService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//添加主任秘钥
func (*leaderKeyApi) Add(r *ghttp.Request) {
	var req *auth_leader_key.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.LeaderkeyService.Add(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//设置主任秘钥状态
func (*leaderKeyApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.LeaderkeyService.SetStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除主任秘钥
func (*leaderKeyApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.LeaderkeyService.Delete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}
