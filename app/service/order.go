package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/order"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/model/order_product"
	"aiyun_cloud_srv/app/model/order_review"
	"aiyun_cloud_srv/app/model/order_upgrade"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"strings"
	"time"
)

var OrderService = new(orderService)

type orderService struct{}

//查询设备信息级订单信息
func (s *orderService) Info(where interface{}) (res *order.Entity, err error) {
	err = order.M.Where(where).Scan(&res)
	return
}

//查询设备信息级订单信息
func (s *orderService) ReviewList(order_number string) (res []*order_review.Entity, err error) {
	err = order_review.M_alias.
		Fields("ore.*,sa.username as operator_name").
		LeftJoin("sys_admin sa", "ore.operator_id = sa.id").
		Where("order_number", order_number).
		Order("ore.id DESC").
		Scan(&res)
	return
}

//查询设备信息级订单信息
func (s *orderService) OrderAndDeviceInfo(device_serial_number string) (o *order.Entity, d *order_device.Entity, err error) {
	err = order_device.M.Where(g.Map{"serial_number": device_serial_number, "status": 1}).OrderDesc("id").Limit(1).Scan(&d)
	if err != nil || d == nil {
		return
	}

	err = order.M.Where("order_number", d.OrderNumber).Scan(&o)

	return
}

//查询订单信息
func (s *orderService) FindOne(where interface{}) (res *order.Entity, err error) {
	err = order.M.Where(where).Limit(1).Scan(&res)
	return
}

func (s *orderService) Create(req *order.Entity, operator_id int) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//创建订单
	now := time.Now().Format("060102150405")
	order_number := fmt.Sprintf("%s%d", now, grand.N(100, 999))
	if _, err := tx.Model(order.Table).Data(g.Map{
		"order_number":     order_number,
		"contract_number":  req.ContractNumber,
		"hospital_id":      req.HospitalId,
		"ukey_code":        req.UkeyCode,
		"sale_name":        req.SaleName,
		"sale_phone":       req.SalePhone,
		"contact":          req.Contact,
		"contact_phone":    req.ContactPhone,
		"receiver":         req.Receiver,
		"receiver_phone":   req.ReceiverPhone,
		"principal":        req.Principal,
		"principal_phone":  req.PrincipalPhone,
		"receiver_region":  req.ReceiverRegion,
		"receiver_address": req.ReceiverAddress,
		"count":            len(req.Devices),
		"maintenance":      req.Maintenance,
		"probation":        req.Probation,
		"operator_id":      operator_id,
		"remark":           req.Remark,
		"status":           0,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	//添加设备
	for k, v := range req.Devices {
		serial_number := fmt.Sprintf("%s%02d", order_number, k+1)
		if _, err := tx.Model(order_device.Table).Data(g.Map{
			"order_number":  order_number,
			"hospital_id":   req.HospitalId,
			"serial_number": serial_number,
			"set_meal":      v.SetMeal,
			"use_type":      v.UseType,
			"remark":        v.Remark,
			"status":        1,
		}).Insert(); err != nil {
			tx.Rollback()
			return err
		}

		//添加产品
		for _, vv := range v.Products {
			for _, vvv := range vv.Children {
				if vvv.MonthCount != 0 {
					if _, err := tx.Model(order_product.Table).Data(g.Map{
						"order_number":  order_number,
						"serial_number": serial_number,
						"set_meal":      v.SetMeal,
						"product_id":    vvv.ProductId,
						"month_count":   vvv.MonthCount,
						"status":        1,
					}).Insert(); err != nil {
						tx.Rollback()
						return err
					}
				}
			}

		}
	}

	return tx.Commit()
}

//设置数据状态
func (s *orderService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = order.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *orderService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = order.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取分页列表
func (s *orderService) Page(req *order.PageReqParams) (total int, list []*order.Entity, err error) {
	M := order.M_alias

	if req.SerialNumber != "" {
		device_list, _ := OrderDeviceService.List(g.Map{"serial_number": req.SerialNumber}, "id")
		order_number_arr := []string{}
		for _, v := range device_list {
			order_number_arr = append(order_number_arr, v.OrderNumber)
		}

		if len(order_number_arr) > 0 {
			M = M.WhereIn("o.order_number", order_number_arr)
		}
	}

	if req.KeyWord != "" {
		M = M.Where("(o.order_number LIKE '%" + req.KeyWord + "%' OR " +
			"o.contract_number LIKE '%" + req.KeyWord + "%' OR " +
			"o.ukey_code LIKE '%" + req.KeyWord + "%' OR " +
			"o.sale_name LIKE '%" + req.KeyWord + "%' OR " +
			"o.sale_phone LIKE '%" + req.KeyWord + "%' OR " +
			"h.name LIKE '%" + req.KeyWord + "%' )")
		//M = M.WhereOrLike("o.order_number", "%"+req.KeyWord+"%")
		//M = M.WhereOrLike("o.contract_number", "%"+req.KeyWord+"%")
		//M = M.WhereOrLike("o.ukey_code", "%"+req.KeyWord+"%")
		//M = M.WhereOrLike("o.sale_name", "%"+req.KeyWord+"%")
		//M = M.WhereOrLike("o.sale_phone", "%"+req.KeyWord+"%")
		//M = M.WhereOrLike("h.name", "%"+req.KeyWord+"%")
	}

	if req.OrderNumber != "" {
		M = M.WhereLike("o.order_number", "%"+req.OrderNumber+"%")
	}

	if req.ContractNumber != "" {
		M = M.WhereLike("o.contract_number", "%"+req.ContractNumber+"%")
	}

	if req.HospitalId > 0 {
		M = M.Where("o.hospital_id", req.HospitalId)
	} else if len(req.RegionIds) == 2 {
		M = M.Where("h.region_id", strings.Join(req.RegionIds, ","))
	}

	if req.Status != -1 {
		M = M.Where("o.status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("o.create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("o.create_at", req.EndTime)
	}

	if len(req.Time) == 2 {
		M = M.WhereGTE("o.update_at", req.Time[0])
		M = M.WhereLTE("o.update_at", req.Time[1])
	}

	if req.Order != "" && req.Sort != "" {
		M = M.Order("o." + req.Order + " " + req.Sort)
	}

	if req.SetMeal > 0 {
		device_list := []*order_device.Entity{}
		err := order_device.M.Where("set_meal", req.SetMeal).Where("status", 1).Scan(&device_list)
		if err != nil {
			g.Log().Error(err)
			err = gerror.New("查询失败")
		}

		order_numbers := []string{}
		for _, v := range device_list {
			order_numbers = append(order_numbers, v.OrderNumber)
		}

		if len(order_numbers) > 0 {
			M = M.Where("o.order_number", order_numbers)
		}
	}

	if req.UseType > 0 {
		device_list := []*order_device.Entity{}
		err := order_device.M.Where("use_type", req.UseType).Where("status", 1).Scan(&device_list)
		if err != nil {
			g.Log().Error(err)
			err = gerror.New("查询失败")
		}

		order_numbers := []string{}
		for _, v := range device_list {
			order_numbers = append(order_numbers, v.OrderNumber)
		}

		if len(order_numbers) > 0 {
			M = M.Where("o.order_number", order_numbers)
		}
	}

	total, err = M.LeftJoin("hospital h", "h.id = o.hospital_id").Group("o.id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("o.*,h.name as hospital_name").LeftJoin("hospital h", "h.id = o.hospital_id").All()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}

	list = make([]*order.Entity, len(data))
	if err = data.Structs(&list); err != nil {
		return
	}

	//获取设备信息
	order_number_arr := []string{}
	for _, v := range list {
		order_number_arr = append(order_number_arr, v.OrderNumber)
	}

	devices, _ := OrderDeviceService.List(g.Map{"order_number": order_number_arr}, "id")
	for _, v := range list {
		for _, v2 := range devices {
			if v.OrderNumber == v2.OrderNumber {
				v.Devices = append(v.Devices, v2)
			}
		}
	}
	return
}

//订单产品补全
func (m *orderService) Parse2ActivateDevice(d *order_device.Entity) (device *order_device.Entity, err error) {

	hospital, err := HospitalService.Info(g.Map{"id": d.HospitalId})
	if err != nil {
		return
	}
	d.HospitalName = hospital.Name
	d.SetMealName = order_product.MealMap[d.SetMeal]

	//旧套餐数据显示
	if d.SetMeal == 1 || d.SetMeal == 2 {
		//初始化模块配置
		//modules, err := initModules()
		//if err != nil {
		//	logging.Error("fail to parse 'conf/modules.json'")
		//	return nil, err
		//}
		//转换模块名称和时间
		//for _, product := range d.Products {
		//	product.DueTimeStr = strconv.Itoa(int(product.DueTime))
		//	info := modules[fmt.Sprintf("%d", product.ProductId+1000)]
		//	product.ProductName = info.Name
		//}

		return nil, nil
	}
	product_list := order_product.AllProductList

	for _, v := range product_list {
		for _, c := range v.Children {
			for _, vv := range d.Products {
				if c.ProductId == vv.ProductId {
					c.MonthCount = vv.MonthCount
					c.DueTime = vv.DueTime
				}
			}
		}
		for _, vv := range d.Products {
			if v.ProductId == vv.ProductId {
				v.MonthCount = vv.MonthCount
				v.DueTime = vv.DueTime
			}
		}
	}

	d.Products = product_list
	device = d
	device.Products = product_list
	return
}

//审核
func (m *orderService) Review(req *order.ReviewReq, operator_id int) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//添加审核记录
	if _, err := tx.Model(order_review.Table).Data(g.Map{
		"order_number": req.OrderNumber,
		"operator_id":  operator_id,
		"status":       req.Status,
		"remark":       req.Remark,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	//修改订单状态
	if _, err := tx.Model(order.Table).
		Where(g.Map{"order_number": req.OrderNumber}).
		Data(g.Map{
			"status": req.Status,
		}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//撤销
func (m *orderService) Revoke(req *order.RevokeReq, o *order.Entity, operator_id int) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	status := 8
	if o.Status == 0 || o.Status == 3 {
		status = 7
	}

	//添加撤销记录
	if _, err := tx.Model(order_review.Table).Data(g.Map{
		"order_number": req.OrderNumber,
		"operator_id":  operator_id,
		"status":       status,
		"remark":       req.Remark,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	//修改订单状态
	if _, err := tx.Model(order.Table).
		Where(g.Map{"order_number": req.OrderNumber}).
		Data(g.Map{
			"status": status,
		}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	//如果是升级的订单 还需要恢复前订单全部数据状态
	if o.PrevOrderNumber != "" {
		//查询当前订单下的所有设备,恢复订单升级前被升级设备的状态
		devices := []*order_device.Entity{}
		err = tx.Model(order_device.Table).Where(g.Map{"order_number": req.OrderNumber, "status": 1}).Scan(&devices)
		if err != nil {
			tx.Rollback()
			return err
		}

		devices_serial_numbers := []string{}
		for _, v := range devices {
			devices_serial_numbers = append(devices_serial_numbers, v.SerialNumber)
		}

		//修改旧设备状态
		if _, err := tx.Model(order_device.Table).
			Where(g.Map{"order_number": o.PrevOrderNumber, "serial_number": devices_serial_numbers}).
			Data(g.Map{
				"status": 1,
			}).Update(); err != nil {
			tx.Rollback()
			return err
		}

		//修改订单状态
		if _, err := tx.Model(order.Table).
			Where(g.Map{"order_number": o.PrevOrderNumber}).
			Data(g.Map{
				"count": gdb.Raw(fmt.Sprintf("count+%d", o.Count)),
			}).Update(); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.Model(order_device.Table).
		Where(g.Map{"order_number": req.OrderNumber}).
		Data(g.Map{
			"status": 0,
		}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

//撤销
func (m *orderService) Deploy(req *order.DeployReq, operator_id int) error {
	//修改订单状态
	_, err := order.M.Where("order_number", req.OrderNumber).Data(g.Map{
		"ukey_code":          req.UkeyCode,
		"express_type":       req.ExpressType,
		"express_no":         req.ExpressNo,
		"deploy_remark":      req.Remark,
		"deploy_operator_id": operator_id,
		"deploy_at":          gtime.Datetime(),
		"status":             9,
	}).Update()
	return err
}

func (s *orderService) Upgrade(req *order.UpgradeReq, devices []*order_device.Entity, o *order.Entity, operator_id int) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	for _, v := range devices {
		//禁用旧设备
		if _, err := tx.Model(order_device.Table).Where(g.Map{
			"order_number":  req.OrderNumber,
			"serial_number": v.SerialNumber,
		}).Data(g.Map{"status": 0}).Update(); err != nil {
			tx.Rollback()
			return err
		}

		//删除旧产品
		//if _, err := tx.Model(order_product.Table).Where(g.Map{
		//	"order_number":  req.OrderNumber,
		//	"serial_number": v.SerialNumber,
		//}).Data(g.Map{"status": 0}).Update(); err != nil {
		//	tx.Rollback()
		//	return err
		//}
		//
		//if _, err := tx.Model(order_device.Table).Data(g.Map{
		//	"order_number":   req.OrderNumber,
		//	"hospital_id":    hospitai_id,
		//	"serial_number":  v.SerialNumber,
		//	"set_meal":       v.SetMeal,
		//	"use_type":       v.UseType,
		//	"remark":         v.Remark,
		//	"machine_brand":  v.MachineBrand,
		//	"machine_number": v.MachineNumber,
		//	"floor":          v.Floor,
		//	"room":           v.Room,
		//	"authorize_time": v.AuthorizeTime,
		//	"activate_time":  v.ActivateTime,
		//	"status":         1,
		//}).Insert(); err != nil {
		//	tx.Rollback()
		//	return err
		//}
		//
		////添加产品
		//for _, vv := range v.Products {
		//	for _, vvv := range vv.Children {
		//		if vvv.MonthCount != 0 {
		//			if _, err := tx.Model(order_product.Table).Data(g.Map{
		//				"order_number":  req.OrderNumber,
		//				"serial_number": v.SerialNumber,
		//				"set_meal":      v.SetMeal,
		//				"product_id":    vvv.ProductId,
		//				"month_count":   vvv.MonthCount,
		//				"due_time":      vvv.DueTime,
		//				"status":        1,
		//			}).Insert(); err != nil {
		//				tx.Rollback()
		//				return err
		//			}
		//		}
		//	}
		//
		//}
	}

	//修改订单审核状态
	//if _, err := tx.Model(order.Table).Where(g.Map{
	//	"order_number": req.OrderNumber,
	//}).Data(g.Map{"status": 0}).Update(); err != nil {
	//	tx.Rollback()
	//	return err
	//}

	//更新被升级订单设备数量
	if _, err := tx.Model(order.Table).Where(g.Map{
		"order_number": req.OrderNumber,
	}).Data(g.Map{"count": gdb.Raw(fmt.Sprintf("count-%d", len(devices)))}).Update(); err != nil {
		tx.Rollback()
		return err
	}

	//if _, err := tx.Model(order_product.Table).Where(g.Map{
	//	"order_number": req.OrderNumber,
	//}).Data(g.Map{"status":0}).Update(); err != nil {
	//	tx.Rollback()
	//	return err
	//}

	//创建订单
	now := time.Now().Format("060102150405")
	order_number := fmt.Sprintf("%s%d", now, grand.N(100, 999))
	if _, err := tx.Model(order.Table).Data(g.Map{
		"order_number":      order_number,
		"prev_order_number": req.OrderNumber,
		"contract_number":   req.ContractNumber,
		"hospital_id":       o.HospitalId,
		"ukey_code":         o.UkeyCode,
		"sale_name":         req.SaleName,
		"sale_phone":        req.SalePhone,
		"contact":           o.Contact,
		"contact_phone":     o.ContactPhone,
		"receiver":          o.Receiver,
		"receiver_phone":    o.ReceiverPhone,
		"principal":         o.Principal,
		"principal_phone":   o.PrincipalPhone,
		"receiver_region":   o.ReceiverRegion,
		"receiver_address":  o.ReceiverAddress,
		"count":             len(devices),
		"maintenance":       o.Maintenance,
		"probation":         o.Probation,
		"operator_id":       operator_id,
		"remark":            o.Remark,
		"status":            0,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	//添加设备
	for _, v := range devices {
		if _, err := tx.Model(order_device.Table).Data(g.Map{
			"order_number":   order_number,
			"hospital_id":    v.HospitalId,
			"serial_number":  v.SerialNumber,
			"set_meal":       v.SetMeal,
			"use_type":       v.UseType,
			"authorize_time": v.AuthorizeTime,
			"machine_brand":  v.MachineBrand,
			"machine_number": v.MachineNumber,
			"floor":          v.Floor,
			"room":           v.Room,
			"remark":         v.Remark,
			"status":         1,
		}).Insert(); err != nil {
			tx.Rollback()
			return err
		}

		//添加产品
		for _, vv := range v.Products {
			for _, vvv := range vv.Children {
				if vvv.MonthCount != 0 {
					if _, err := tx.Model(order_product.Table).Data(g.Map{
						"order_number":  order_number,
						"serial_number": v.SerialNumber,
						"set_meal":      v.SetMeal,
						"due_time":      -1,
						"product_id":    vvv.ProductId,
						"month_count":   vvv.MonthCount,
						"status":        1,
					}).Insert(); err != nil {
						tx.Rollback()
						return err
					}
				}
			}

		}
	}

	//添加升级记录
	if _, err := tx.Model(order_upgrade.Table).Data(g.Map{
		"order_number":    req.OrderNumber,
		"contract_number": req.ContractNumber,
		"sale_name":       req.SaleName,
		"sale_phone":      req.SalePhone,
		"operator_id":     operator_id,
		"remaek":          req.Remark,
		"status":          1,
	}).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
