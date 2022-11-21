package Api

import (
	"aiyun_cloud_srv/app/model/auth_client"
	"aiyun_cloud_srv/app/model/auth_ukey"
	"aiyun_cloud_srv/app/model/licence"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"aiyun_cloud_srv/library/utils/rsa_crypt"
	"encoding/base64"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Authorize = authorizeApi{}

type authorizeApi struct{}

// 检测服务端授权状态
func (*authorizeApi) StatusSrv(r *ghttp.Request) {
	var req *auth_ukey.QueryUkeyStatus

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	info, err := service.AuthUkeyService.FindOne(g.Map{
		"serial_number": req.SerialNumber,
		"status":        1,
	})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if info == nil {
		response.Json(r, 7, "Ukey未注册")
	}

	//验证签名
	pubKey, err := rsa_crypt.LoadPublicKeyBase64(info.PublicKey)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if !rsa_crypt.Verify(pubKey, utils.EncodeSHA256(req.Uuid+req.SerialNumber, ""), req.Signature) {
		response.Json(r, 9, "验签失败，请检查该U-Key是否注册")
	}

	if info.AuthTimes-info.UsedTimes <= 0 {
		response.Json(r, 6, "此U-Key的授权次数已用尽")
	}

	response.Success(r, info)
}

// 查询客户端授权状态，如果已经授权 则返回授权信息给客户端
func (*authorizeApi) Status(r *ghttp.Request) {
	var req *licence.ClientAuthStatus

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	client_info, err := service.AuthClientService.FindByUUID(req.Uuid, 1)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if client_info == nil {
		response.Json(r, 0, "未授权", g.Map{"state": 0})
	}

	// 获取授权u_key信息
	ukey_info, err := service.AuthUkeyService.FindOne(g.Map{"serial_number": client_info.UkeySerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if ukey_info == nil {
		response.Json(r, 7, "Ukey未注册")
	}

	licence_info, err := service.AuthLicenceService.AuthExist(req.Uuid, 1)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if licence_info == nil {
		response.Json(r, 0, "未授权", g.Map{"state": 0})
	}

	response.Json(r, 0, "已授权", g.Map{
		"state":       1,
		"licence":     licence_info.Licence,
		"licence_key": ukey_info.PublicKey,
	})
}

//检查Ukey是否存在, 并返回ukey_code
func (*authorizeApi) UkeyInfo(r *ghttp.Request) {
	ukey_serial_number := r.GetQueryString("ukey_serial_number")

	ukey_info, err := service.AuthUkeyService.FindOne(g.Map{"serial_number": ukey_serial_number, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if ukey_info != nil {
		response.Success(r, g.Map{"ukey_code": ukey_info.Code})
	}

	response.Json(r, 7, "U-Key未注册")
}

//查询授权订单
func (*authorizeApi) QueryOrder(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")
	if serial_number == "" {
		response.Error(r, "参数错误,设备系列号必填")
	}

	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(serial_number)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if device == nil {
		response.Error(r, "设备不存在")
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	if device.ActivateTime > 0 {
		response.Error(r, "此设备已被激活")
	}

	if order.Status < 9 && order.Status != 1 {
		response.Error(r, "订单状态异常，未走完审批流程")
	}

	//if order.Status == 1{
	//	response.Error(r, "订单状态异常，此订单已完成部署")
	//}

	licence, err := service.AuthLicenceService.FindOne(g.Map{"device_serial_number": serial_number})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence != nil {
		if licence.Status == 0 {
			response.Error(r, "此设备已禁用")
		}
	}

	products, err := service.OrderDeviceService.ProductList(g.Map{"order_number": order.OrderNumber, "serial_number": serial_number, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	device.Products = products
	//计算产品模块到期时间
	//todo

	device, err = service.OrderService.Parse2ActivateDevice(device)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, device)
}

//创建服务端授权证书
func (*authorizeApi) CreateSrvLicence(r *ghttp.Request) {
	var req *auth_client.Entity

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//检查医院
	if req.HospitalId <= 0 {
		response.Error(r, "被授权单位必填")
	}

	hospital, err := service.HospitalService.Info(g.Map{"id": req.HospitalId, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if hospital == nil {
		response.Error(r, "单位不存在或被禁止")
	}

	req.HospitalId = hospital.Id

	// 获取授权u_key信息
	ukey_info, err := service.AuthUkeyService.FindOne(g.Map{
		"serial_number": req.UkeySerialNumber,
		"status":        1,
	})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if ukey_info == nil {
		response.Json(r, 7, "Ukey未注册")
	}

	//验证签名
	pubKey, err := rsa_crypt.LoadPublicKeyBase64(ukey_info.PublicKey)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if !rsa_crypt.Verify(pubKey, utils.EncodeSHA256(req.Uuid+req.UkeySerialNumber, ""), req.Signature) {
		response.Json(r, 9, "验签失败，请检查该U-Key是否已被创建")
	}

	if ukey_info.AuthTimes-ukey_info.UsedTimes <= 0 {
		response.Json(r, 6, "此U-Key的授权次数已用尽")
	}

	// 查询授权证书是否存在，如果不存在则创建，如果存在则直接返回证书信息
	licence_info, err := service.AuthLicenceService.AuthExist(req.Uuid, req.Role)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if licence_info != nil {
		response.Success(r, g.Map{
			"author_number": licence_info.AuthorNumber,
			"licence":       licence_info.Licence,
			"licence_key":   licence_info.PublicKey,
			"ukey_code":     ukey_info.Code,
			"ukey_crypt":    rsa_crypt.Crypt(pubKey, ukey_info.SerialNumber),
		})
	}

	// 创建证书
	licence_info, err = service.AuthLicenceService.Create(req, ukey_info.Code, nil)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence_info == nil {
		response.Error(r, "创建授权失败")
	}

	response.Success(r, g.Map{
		"author_number": licence_info.AuthorNumber,
		"licence":       licence_info.Licence,
		"licence_key":   licence_info.PublicKey,
		"ukey_code":     ukey_info.Code,
		"ukey_crypt":    rsa_crypt.Crypt(pubKey, ukey_info.SerialNumber),
	})
}

// 创建客户端授权证书
func (*authorizeApi) CreateLicence(r *ghttp.Request) {
	var req *auth_client.Entity

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//客户端授权
	req.Role = 1

	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(req.SerialNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if device == nil {
		response.Error(r, "设备不存在")
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	if device.AuthorizeTime > 0 {
		response.Error(r, "此设备已被授权，禁止再次授权")
	}

	if order.Status < 9 && order.Status != 1 {
		response.Error(r, "订单状态异常，未走完审批流程")
	}

	// 获取授权u_key信息
	ukey_info, err := service.AuthUkeyService.FindOne(g.Map{"serial_number": req.UkeySerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if ukey_info == nil {
		response.Error(r, "Ukey未注册")
	}

	//验证签名
	pubKey, err := rsa_crypt.LoadPublicKeyBase64(ukey_info.PublicKey)
	if err != nil {
		response.ErrorSys(r, err)
	}

	msg := utils.GenSignMsg(auth_client.SignData{
		Uuid:             req.Uuid,
		SerialNumber:     req.SerialNumber,
		UkeySerailNumber: req.UkeySerialNumber,
		CpuID:            req.CpuID,
		CpuName:          req.CpuName,
		BaseboardID:      req.BaseboardID,
	})

	if !rsa_crypt.Verify(pubKey, utils.EncodeSHA256(msg, ""), req.Signature) {
		response.Error(r, "验签失败，请检查该U-Key是否已被创建")
	}

	if ukey_info.AuthTimes-ukey_info.UsedTimes <= 0 {
		response.Error(r, "此U-Key的授权次数已用尽")
	}

	// 查询授权证书是否存在，如果不存在则直接创建，如果存在，直接返回证书信息
	licence_info, err := service.AuthLicenceService.AuthExist(req.Uuid, req.Role)
	if err != nil {
		response.ErrorSys(r, err)
	}

	if licence_info != nil {
		response.Success(r, g.Map{
			"licence":     licence_info.Licence,
			"licence_key": licence_info.PublicKey,
			"state":       1,
		})
	}

	req.HospitalId = order.HospitalId

	// 创建证书
	licence_info, err = service.AuthLicenceService.Create(req, ukey_info.Code, order)
	if err != nil {
		response.ErrorDb(r, err)
	}
	if licence_info == nil {
		response.Error(r, "创建授权失败")
	}

	response.Success(r, g.Map{
		"licence":     licence_info.Licence,
		"licence_key": licence_info.PublicKey,
		"state":       1,
	})
}

// 激活/升级状态
func (*authorizeApi) ActivateStatus(r *ghttp.Request) {
	var req *licence.ActivateStatusReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询证书是否存在
	licence_info, err := service.AuthLicenceService.FindOne(
		g.Map{"author_number": req.AuthorNumber, "device_serial_number": req.SerialNumber, "status": 1})
	if err != nil {
		response.ErrorSys(r, err)
	}
	if licence_info == nil {
		response.Json(r, 8, "客户端未授权")
	}

	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(req.SerialNumber)
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
		response.Error(r, "此设备未授权")
	}

	//if order.Status != 9 {
	//	response.Error(r, "订单状态异常，未走完审批流程")
	//}

	// 未升级、已升级的订单，从未激活过， 返回状态 0：未激活,1：已激活
	// 已升级的订单，之前激活过， 返回状态 1：已激活,2：未激活

	if device.ActivateTime > 0 {
		products, err := service.OrderDeviceService.ProductList(g.Map{"order_number": order.OrderNumber, "serial_number": req.SerialNumber, "status": 1})
		if err != nil {
			response.ErrorDb(r, err)
		}
		device.Products = products
		//计算产品模块到期时间
		//todo

		device, err = service.OrderService.Parse2ActivateDevice(device)
		if err != nil {
			response.ErrorDb(r, err)
		}

		b, err := json.Marshal(device)
		if err != nil {
			response.ErrorSys(r, err)
		}

		response.Json(r, 0, "设备已激活", g.Map{"state": 1, "device_info": base64.StdEncoding.EncodeToString(b)})
	}

	if order.PrevOrderNumber != "" {
		d, err := service.OrderDeviceService.FindOne(g.Map{"order_number": order.PrevOrderNumber, "serial_number": req.SerialNumber})
		if err != nil {
			response.ErrorDb(r, err)
		}
		if d.AuthorizeTime > 0 {
			products, err := service.OrderDeviceService.ProductList(g.Map{"order_number": order.OrderNumber, "serial_number": req.SerialNumber, "status": 1})
			if err != nil {
				response.ErrorDb(r, err)
			}
			device.Products = products
			//计算产品模块到期时间
			//todo

			device, err = service.OrderService.Parse2ActivateDevice(device)
			if err != nil {
				response.ErrorDb(r, err)
			}

			b, err := json.Marshal(device)
			if err != nil {
				response.ErrorSys(r, err)
			}

			response.Json(r, 0, "设备未激活", g.Map{"state": 2, "device_info": base64.StdEncoding.EncodeToString(b)})
		}
	}

	response.Json(r, 0, "设备未激活", g.Map{"state": 0, "device_info": nil})
}

// 设备激活
func (*authorizeApi) Activate(r *ghttp.Request) {
	var req *licence.ActivateReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询证书是否存在
	licence_info, err := service.AuthLicenceService.FindOne(g.Map{"author_number": req.AuthorNumber, "status": 1})
	if err != nil {
		response.ErrorSys(r, err)
	}
	if licence_info == nil {
		response.Json(r, 8, "客户端未授权")
	}

	//查询订单
	order, device, err := service.OrderService.OrderAndDeviceInfo(req.SerialNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if device == nil {
		response.Error(r, "设备不存在")
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	is_upgrade := true
	if device.MachineBrand == "" { //激活
		is_upgrade = false
		if req.MachineBrand == "" {
			response.Error(r, "超声设备品牌信息必填")
		}

		if req.MachineNumber == "" {
			response.Error(r, "超声设备号必填")
		}

		if req.Floor == "" {
			response.Error(r, "楼层信息必填")
		}

		if req.Room == "" {
			response.Error(r, "房号信息必填")
		}
	}

	if device.AuthorizeTime <= 0 {
		response.Error(r, "此设备未授权，禁止激活")
	}

	if order.Status < 9 && order.Status != 1 {
		response.Error(r, "订单状态异常，未走完审批流程")
	}

	if device.ActivateTime > 0 {
		response.Error(r, "设备已被激活，禁止重复操作")
	}

	//计算模块到期时间
	//todo

	//激活操作
	if err = service.OrderDeviceService.Activate(req, is_upgrade, device); err != nil {
		response.ErrorDb(r, err)
	}

	products, err := service.OrderDeviceService.ProductList(g.Map{"order_number": order.OrderNumber, "serial_number": req.SerialNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}
	device.Products = products
	//计算产品模块到期时间
	//todo

	device, err = service.OrderService.Parse2ActivateDevice(device)
	if err != nil {
		response.ErrorDb(r, err)
	}

	b, err := json.Marshal(device)
	if err != nil {
		response.ErrorSys(r, err)
	}

	response.Success(r, g.Map{"state": 1, "device_info": base64.StdEncoding.EncodeToString(b)})
}

// 设备激活(测试使用)
func (*authorizeApi) CancelActivate(r *ghttp.Request) {
	serial_number := r.GetQueryString("serial_number")
	if serial_number == "" {
		response.Error(r, "设备序列号必填")
	}

	_, err := order_device.M.Where("serial_number", serial_number).Data(g.Map{"activate_time": 0}).Update()
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r)
}
