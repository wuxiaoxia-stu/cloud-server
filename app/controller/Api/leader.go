package Api

import (
	"aiyun_cloud_srv/app/model/auth_leader_key"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var LeaderKey = leaderKeyApi{}

type leaderKeyApi struct{}

// 查询秘钥基本信息
func (*leaderKeyApi) Info(r *ghttp.Request) {
	serial_number := r.GetQueryString("leader_serial_number")

	leader_info, err := service.LeaderkeyService.FindOne(g.Map{"serial_number": serial_number, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_info != nil {
		response.Success(r, g.Map{"leader_key_code": leader_info.Code})
	}

	response.Json(r, 7, "主任秘钥未注册")
}

// 绑定主任秘钥
func (*leaderKeyApi) Bind(r *ghttp.Request) {
	var req *auth_leader_key.BindLeaderReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询客户端证书是否存在
	client_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.AuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if client_licence == nil {
		response.Error(r, "客户端未授权")
	}

	//查询服务端证书是否存在
	service_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ServerAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if service_licence == nil {
		response.Error(r, "服务端未授权")
	}

	////验签数据
	//service_pubKey, err := rsa_crypt.LoadPublicKeyBase64(service_licence.PublicKey)
	//if err != nil {
	//	response.ErrorSys(r, err)
	//}
	//
	//if !rsa_crypt.Verify(service_pubKey, req.ServerAuthorNumber, req.Signature) {
	//	response.Json(r, 9, "验签失败，服务端授权异常")
	//}

	//验证主任秘钥
	leader_key_info, err := service.LeaderkeyService.FindOne(g.Map{"serial_number": req.LeaderSerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_key_info == nil {
		response.Error(r, "主任秘钥未注册")
	}

	//检查是否被绑定
	leader_bind_info, err := service.LeaderkeyService.GetBindInfo(g.Map{"auth_leader_id": leader_key_info.Id, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_bind_info != nil {
		response.Error(r, "此主任秘钥已经被绑定")
	}

	//检查绑定的医院是否匹配
	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(client_licence.DeviceSerialNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if device == nil {
		response.Error(r, "设备不存在")
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	if device.AuthorizeTime <= 0 {
		response.Error(r, "此设备未授权，禁止此操作")
	}

	if !(order.Status == 9 || order.Status == 1) {
		response.Error(r, "订单状态异常，未走完审批流程")
	}

	if leader_key_info.HospitalId != device.HospitalId {
		response.Error(r, "授权单位不匹配")
	}

	//修改主任秘钥授权次数
	if err := service.LeaderkeyService.Bind(req, leader_key_info, device.SerialNumber); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, leader_key_info)
}

// 重新绑定主任秘钥
func (*leaderKeyApi) Rebind(r *ghttp.Request) {
	var req *auth_leader_key.BindLeaderReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询客户端证书是否存在
	client_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.AuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if client_licence == nil {
		response.Error(r, "客户端未授权")
	}

	//查询服务端证书是否存在
	service_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ServerAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if service_licence == nil {
		response.Error(r, "服务端未授权")
	}

	//验证主任秘钥
	leader_key_info, err := service.LeaderkeyService.FindOne(g.Map{"serial_number": req.LeaderSerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_key_info == nil {
		response.Error(r, "主任秘钥未注册")
	}

	//检查是否被绑定
	leader_bind_info, err := service.LeaderkeyService.GetBindInfo(g.Map{"auth_leader_id": leader_key_info.Id, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_bind_info == nil {
		response.Error(r, "此主任秘钥未被绑定，此操作无效")
	}

	//检查绑定的医院是否匹配
	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(client_licence.DeviceSerialNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if device == nil {
		response.Error(r, "设备不存在")
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	if device.AuthorizeTime <= 0 {
		response.Error(r, "此设备未授权，禁止此操作")
	}

	if !(order.Status == 9 || order.Status == 1) {
		response.Error(r, "订单状态异常，未走完审批流程")
	}

	if leader_key_info.HospitalId != device.HospitalId {
		response.Error(r, "授权单位不匹配")
	}

	//修改主任秘钥授权次数
	if err := service.LeaderkeyService.Rebind(req, leader_key_info, device.SerialNumber); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}

// 重新绑定主任秘钥
func (*leaderKeyApi) Unbind(r *ghttp.Request) {
	var req *auth_leader_key.UnbindLeaderReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询服务端证书是否存在
	service_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ServerAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if service_licence == nil {
		response.Error(r, "服务端未授权")
	}

	//验证主任秘钥
	leader_key_info, err := service.LeaderkeyService.FindOne(g.Map{"serial_number": req.LeaderSerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_key_info == nil {
		response.Error(r, "主任秘钥未注册")
	}

	//检查是否被绑定
	leader_bind_info, err := service.LeaderkeyService.GetBindInfo(g.Map{"auth_leader_id": leader_key_info.Id, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if leader_bind_info == nil {
		response.Error(r, "此主任秘钥未被绑定，此操作无效")
	}

	//修改主任秘钥授权次数
	if err := service.LeaderkeyService.Unbind(leader_key_info.Id); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
