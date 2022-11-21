package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/order"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/model/order_product"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Order = orderApi{}

type orderApi struct{}

// 获取订单列表
func (*orderApi) List(r *ghttp.Request) {
	var req *order.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	total, list, err := service.OrderService.Page(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, g.Map{
		"list":  list,
		"total": total,
	})
}

//创建订单
func (*orderApi) Create(r *ghttp.Request) {
	var req *order.Entity

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	count := len(req.Devices)
	if count <= 0 || count >= 100 {
		response.Error(r, "设备数量异常，可添加设备数为1到99之间")
	}

	//校验产品合法性
	for index, device := range req.Devices {
		if !(device.SetMeal >= 3 && device.SetMeal <= 9) {
			response.Error(r, fmt.Sprintf("第%d台设备套餐信息异常", index+1))
		}

		if !utils.InArray(device.UseType, []int{1, 2, 3, 4}) {
			response.Error(r, fmt.Sprintf("第%d台设备设备用途信息异常", index+1))
		}

		if len(device.Products) <= 0 {
			response.Error(r, fmt.Sprintf("第%d台设备产品信息异常", index+1))
		}

		//产品使用日期，和id是否合法
		valid_product_count := 0
		for _, product := range device.Products {
			for _, product_child := range product.Children {
				if !utils.InArray(product_child.ProductId, order_product.AllMealOptionMap) {
					response.Error(r, fmt.Sprintf("第%d台设备产品配置不存在", index+1))
				}

				if product_child.MonthCount != 0 {
					valid_product_count++
				}
			}
		}
		if valid_product_count == 0 {
			response.Error(r, fmt.Sprintf("第%d台设备产品信息异常", index+1))
		}
	}

	err := service.OrderService.Create(req, r.GetCtxVar("uid").Int())
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "添加成功")
}

//删除订单数据
func (*orderApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	_, err := service.OrderService.Delete(req)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//获取配置信息
func (*orderApi) MealOptions(r *ghttp.Request) {
	response.Success(r, g.Map{
		"meal_options":     order_product.MealMap,
		"product_options":  order_product.AllProductList,
		"meal_def_options": order_product.MealOptionDefMap,
	})
}

//订单详情
func (*orderApi) Details(r *ghttp.Request) {
	order_number := r.GetQueryString("order_number")
	order, err := service.OrderService.Info(g.Map{"order_number": order_number})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if order == nil {
		response.Error(r, "订单不存在")
	}

	//查询医院和地区信息
	hospital, err := service.HospitalService.Info(g.Map{"id": order.HospitalId})
	if err != nil {
		response.ErrorDb(r, err)
	}
	hospital_region, err := service.RegionService.GetNameById(hospital.RegionId)
	if err != nil {
		response.ErrorDb(r, err)
	}
	hospital.RegionName = hospital_region
	order.Hospital = hospital

	if order.ReceiverRegion != "" {
		receiver_region, err := service.RegionService.GetNameById(order.ReceiverRegion)
		if err != nil {
			response.ErrorDb(r, err)
		}
		order.ReceiverRegionName = receiver_region
	}

	//获取设备列表
	devices, err := service.OrderDeviceService.List(g.Map{"order_number": order.OrderNumber}, "serial_number")
	if err != nil {
		response.ErrorDb(r, err)
	}
	//获取产品列表
	products, err := service.OrderDeviceService.ProductList(g.Map{"order_number": order.OrderNumber, "status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	for _, v := range devices {
		for _, v2 := range products {
			if v.SerialNumber == v2.SerialNumber {
				v.Products = append(v.Products, v2)
			}
		}
	}

	product_all := order_product.AllProductList
	for k, v := range devices {
		data, _ := json.Marshal(product_all)
		var product_list []*order_product.Entity
		err := json.Unmarshal(data, &product_list)
		if err != nil {
			response.ErrorSys(r, err)
		}
		for k2, v2 := range product_list {
			for c2, c := range v2.Children {
				product_list[k2].Children[c2].MonthCount = 0
				for _, vv := range v.Products {
					if c.ProductId == vv.ProductId {
						c.MonthCount = vv.MonthCount
						c.DueTime = vv.DueTime
					}
				}
			}
		}
		devices[k].Products = product_list
	}

	order.Devices = devices

	////审核数据
	review_list, err := service.OrderService.ReviewList(order.OrderNumber)
	if err != nil {
		response.ErrorDb(r, err)
	}

	order.Reviews = review_list
	response.Success(r, order)
}

// 订单审核
func (*orderApi) Review(r *ghttp.Request) {
	var req *order.ReviewReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询订单状态
	order, err := service.OrderService.Info(g.Map{"order_number": req.OrderNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if order == nil {
		response.Error(r, "订单不存在")
	}
	if !(order.Status == 0 || order.Status == 2) {
		response.Error(r, "订单状态异常,此状态禁审核操作")
	}
	if order.Status == 0 && !(req.Status == 2 || req.Status == 3) {
		response.Error(r, "参数错误")
	}
	if order.Status == 2 && !(req.Status == 5 || req.Status == 6) {
		response.Error(r, "参数错误")
	}

	if err := service.OrderService.Review(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "审核完成")
}

//撤销
func (*orderApi) Revoke(r *ghttp.Request) {
	var req *order.RevokeReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询订单状态
	order, err := service.OrderService.Info(g.Map{"order_number": req.OrderNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if order == nil {
		response.Error(r, "订单不存在")
	}
	if order.Status == 1 || order.Status == 9 {
		response.Error(r, "此状态禁止撤销")
	}

	if err := service.OrderService.Revoke(req, order, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "审核完成")
}

//订单升级
func (*orderApi) Deploy(r *ghttp.Request) {
	var req *order.DeployReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//查询订单状态
	order, err := service.OrderService.Info(g.Map{"order_number": req.OrderNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if order == nil {
		response.Error(r, "订单不存在")
	}

	if order.Status != 5 {
		response.Error(r, "此状态禁止此操作")
	}

	if err := service.OrderService.Deploy(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "审核完成")
}

//订单升级
func (*orderApi) Upgrade(r *ghttp.Request) {
	var req *order.UpgradeReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}
	count := len(req.Devices)
	if count <= 0 || count >= 100 {
		response.Error(r, "设备数量异常，可添加设备数为1到99之间")
	}

	//校验产品合法性
	for index, device := range req.Devices {
		if !(device.SetMeal >= 3 && device.SetMeal <= 9) {
			response.Error(r, fmt.Sprintf("第%d台设备套餐信息异常", index+1))
		}

		if !utils.InArray(device.UseType, []int{1, 2, 3, 4}) {
			response.Error(r, fmt.Sprintf("第%d台设备设备用途信息异常", index+1))
		}
		if len(device.Products) <= 0 {
			response.Error(r, fmt.Sprintf("第%d台设备产品信息异常", index+1))
		}

		//产品使用日期，和id是否合法
		valid_product_count := 0
		for _, product := range device.Products {
			for _, product_child := range product.Children {
				if !utils.InArray(product_child.ProductId, order_product.AllMealOptionMap) {
					response.Error(r, fmt.Sprintf("第%d台设备产品配置不存在", index+1))
				}

				if product_child.MonthCount != 0 {
					valid_product_count++
				}
			}
		}
		if valid_product_count == 0 {
			response.Error(r, fmt.Sprintf("第%d台设备产品信息异常", index+1))
		}
	}

	//查询订单状态
	order, err := service.OrderService.Info(g.Map{"order_number": req.OrderNumber})
	if err != nil {
		response.ErrorDb(r, err)
	}
	if order == nil {
		response.Error(r, "订单不存在")
	}

	if !(order.Status == 9 || order.Status == 1) {
		response.Error(r, "此状态禁止升级")
	}

	//查询当前设备列表和提交的设备信息比对，剔除套餐未变动的设备
	devices, err := service.OrderDeviceService.List(g.Map{"order_number": req.OrderNumber, "status": 1}, "id")
	if err != nil {
		response.ErrorDb(r, err)
	}

	//需要升级的设备
	upgrade_device_list := []*order_device.Entity{}
	for _, v := range devices {
		for _, v2 := range req.Devices {
			if v.SerialNumber == v2.SerialNumber {
				if v2.SetMeal != 9 && v.SetMeal < v2.SetMeal {
					v.SetMeal = v2.SetMeal
					v.UseType = v2.UseType
					v.Products = v2.Products
					upgrade_device_list = append(upgrade_device_list, v)
				}
				if v.SetMeal == 9 && v2.SetMeal != 9 {
					v.SetMeal = v2.SetMeal
					v.UseType = v2.UseType
					v.Products = v2.Products
					upgrade_device_list = append(upgrade_device_list, v)
				}
			}
		}
	}

	if len(upgrade_device_list) <= 0 {
		response.Error(r, "未检测到升级的设备，请至少选择一台设备升级")
	}

	if err := service.OrderService.Upgrade(req, upgrade_device_list, order, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "订单升级成功")
}
