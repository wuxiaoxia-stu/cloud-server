package order_device

import (
	"aiyun_cloud_srv/app/model/order_product"
	"github.com/gogf/gf/frame/g"
)

// Entity is the golang structure for table order_device
type Entity struct {
	Id            int                     `orm:"id,primary,table_comment:'设备管理'" json:"id"`
	OrderNumber   string                  `orm:"order_number,size:15,not null,comment:'订单号'" json:"order_number"`
	HospitalId    int                     `orm:"hospital_id,size:4,not null,comment:'单位id'"`
	SetMeal       int                     `orm:"set_meal,size:2,comment:'套餐'" json:"set_meal"`
	DefaultMonth  int                     `orm:"default_month,size:2,comment:'标配模块默认月份'" json:"default_month"`
	SerialNumber  string                  `orm:"serial_number,size:18,not null,comment:'设备sn码'" json:"serial_number"`
	DeviceType    int                     `orm:"device_type,size:2,comment:'设备类型 0新设备 1已有设备'" json:"device_type"` // 设备类型 0新设备 1已有设备
	Typee         int                     `orm:"typee,size:2,comment:'授权类型 0新设备授权 1升级续费'" json:"typee"`           // 授权类型 0新设备授权 1升级续费
	UseType       int                     `orm:"use_type,size:2,not null,default:4,comment:'设备用途'" json:"use_type"`
	MachineBrand  string                  `orm:"machine_brand,size:100,comment:'超声机品牌'" json:"machine_brand"`
	MachineNumber string                  `orm:"machine_number,size:20,comment:'超声机'" json:"machine_number"`
	Floor         string                  `orm:"floor,size:20,comment:'楼层'" json:"floor"`
	Room          string                  `orm:"room,size:20,comment:'房号'" json:"room" `
	AuthorizeTime int                     `orm:"authorize_time,size:4,comment:'授权时间'" json:"authorize_time"`
	ActivateTime  int                     `orm:"activate_time,size:4,comment:'激活时间'" json:"activate_time"`
	Remark        string                  `orm:"remark,size:500,comment:'备注'" json:"remark"`
	Status        int                     `orm:"status,size:2,comment:'状态'" json:"status"`
	Products      []*order_product.Entity `json:"products"`
	HospitalName  string                  `json:"hospital_name"`
	SetMealName   string                  `json:"set_meal_name"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "order_device"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "order_device od"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
