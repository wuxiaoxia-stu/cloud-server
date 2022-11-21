package order_upgrade

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table order_review
type Entity struct {
	Id             int         `orm:"id,primary,table_comment:'设备管理'" json:"id"`
	OrderNumber    string      `orm:"order_number,size:15,not null,comment:'订单号'"`
	ContractNumber string      `orm:"contract_number,size:15,not null,comment:'合同号'"`
	SaleName       string      `orm:"sale_name,size:20,not null,comment:'销售名称'"`
	SalePhone      string      `orm:"sale_phone,size:20,not null,comment:'销售电话'"`
	OperatorId     int         `orm:"operator_id,not null,comment:'操作人'"`
	Remark         string      `orm:"remark,size:500,comment:'备注'" json:"remark"`
	UpdateAt       *gtime.Time `orm:"update_at" json:"update_at"`
	Status         int         `orm:"status,size:2,comment:'状态'" json:"status"`
}

var (
	// Table is the table name of order_review.
	Table       = "order_upgrade"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "order_upgrade op"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
