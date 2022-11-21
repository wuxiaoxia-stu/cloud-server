package Api

import (
	"aiyun_cloud_srv/app/model/auth_licence"
	"aiyun_cloud_srv/app/model/licence"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/rsa_crypt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

var Pair = pairApi{}

type pairApi struct{}

// 处理配对
func (*pairApi) Apply(r *ghttp.Request) {
	var req *licence.PairReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询授权记录是否存在
	service_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ServerAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if service_licence == nil {
		response.Error(r, "服务端未授权")
	}

	client_licence, err := service.AuthLicenceService.FindOne(g.Map{"device_serial_number": req.ClientSerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if client_licence == nil {
		response.Error(r, "客户端未授权")
	}

	if client_licence.AuthorNumber != req.ClientAuthorNumber {
		response.Error(r, "客户端授权异常，授权码不匹配")
	}

	if client_licence.HospitalId != service_licence.HospitalId {
		response.Error(r, "客户端授权单位与服务端授权单位不匹配")
	}

	//解密数据比较
	service_priKey, err := rsa_crypt.LoadPrivateKeyBase64(service_licence.PrivateKey)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if rsa_crypt.DeCrypt(service_priKey, req.Signature) != utils.GenSignMsg(req) {
		response.Json(r, 9, "数据异常")
	}

	//更新授权记录表，把服务端授权id绑定到客户端
	if _, err = auth_licence.M.Where("id", client_licence.Id).Data(g.Map{
		"licence_id": service_licence.Id,
		"pair_at":    gtime.Datetime(),
	}).Update(); err != nil {
		response.ErrorDb(r, err)
	}

	licence_client_info, err := service.AuthLicenceService.InfoAll(client_licence.DeviceSerialNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"uuid":          licence_client_info.ClientInfo.Uuid,
		"ukey_code":     licence_client_info.UkeyCode,
		"author_number": licence_client_info.AuthorNumber,
		"device_number": client_licence.DeviceSerialNumber,
		"public_key":    licence_client_info.PublicKey,
	})
}

// 解除配对
func (*pairApi) Break(r *ghttp.Request) {
	var req *licence.BreakPairReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询授权记录是否存在
	service_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ServerAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if service_licence == nil {
		response.Error(r, "服务端未授权")
	}

	client_licence, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.ClientAuthorNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if client_licence == nil {
		response.Error(r, "客户端未授权")
	}

	//解密数据比较
	service_priKey, err := rsa_crypt.LoadPrivateKeyBase64(service_licence.PrivateKey)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if rsa_crypt.DeCrypt(service_priKey, req.Signature) != utils.GenSignMsg(req) {
		response.Json(r, 9, "数据异常")
	}

	//更新授权记录表，把服务端授权id绑定到客户端
	if _, err = auth_licence.M.Where("id", client_licence.Id).Data(g.Map{"licence_id": 0}).Update(); err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
