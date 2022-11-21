package service

import (
	"aiyun_cloud_srv/app/model/licence"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/model/order_product"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"time"
)

var OrderDeviceService = new(orderDeviceService)

type orderDeviceService struct{}

//查询设备信息
func (s *orderDeviceService) FindOne(where interface{}) (res *order_device.Entity, err error) {
	err = order_device.M.Where(where).Limit(1).Scan(&res)
	return
}

func (s *orderDeviceService) List(where interface{}, order string) (res []*order_device.Entity, err error) {
	err = order_device.M.Where(where).Order(order).Scan(&res)
	return
}

func (s *orderDeviceService) ProductList(where interface{}) (res []*order_product.Entity, err error) {
	err = order_product.M.Where(where).Scan(&res)
	return
}

func (s *orderDeviceService) Activate(req *licence.ActivateReq, is_upgrade bool, device *order_device.Entity) error {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	data := g.Map{
		"activate_time": time.Now().Unix(),
	}

	if !is_upgrade {
		data["machine_brand"] = req.MachineBrand
		data["machine_number"] = req.MachineNumber
		data["floor"] = req.Floor
		data["room"] = req.Room
	}

	//设置激活时间
	if _, err = tx.Model(order_device.Table).Where(g.Map{
		"serial_number": req.SerialNumber,
		"status":        1,
	}).Data(data).Update(); err != nil {
		tx.Rollback()
		return err
	}

	//设置产品到期时间 due_time
	if device.SetMeal == 9 {
		products := []*order_product.Entity{}
		if err = tx.Model(order_product.Table).Where(g.Map{
			"serial_number": req.SerialNumber,
			"order_number":  device.OrderNumber,
			"status":        1,
		}).Scan(&products); err != nil {
			tx.Rollback()
			return err
		}
		if len(products) <= 0 {
			tx.Rollback()
			return fmt.Errorf("设备模块配置异常")
		}

		for _, v := range products {
			if _, err = tx.Model(order_product.Table).Where("id", v.Id).Data(g.Map{
				"due_time": time.Now().AddDate(0, v.MonthCount, 0).Unix(),
			}).Update(); err != nil {
				tx.Rollback()
				return err
			}
		}

	} else {
		if _, err := tx.Model(order_product.Table).Where(g.Map{
			"serial_number": req.SerialNumber,
			"order_number":  device.OrderNumber,
			"status":        1,
		}).Data(g.Map{"due_time": -1}).Update(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
