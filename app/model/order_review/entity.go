package order_review

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table order_review
type Entity struct {
	Id int `orm:"id,primary,table_comment:'设备管理'" json:"id"`
	//Type        int         `orm:"type,size:2,comment:'类型：1：财务审核'，2：总经理审核，3：撤销"`
	OrderNumber  string      `orm:"order_number,size:15,not null,comment:'订单号'"`
	OperatorId   int         `orm:"operator_id,not null,comment:'操作人'"`
	Remark       string      `orm:"remark,size:500,comment:'备注'" json:"remark"`
	CreateAt     *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at" json:"update_at"`
	Status       int         `orm:"status,size:2,comment:'状态'" json:"status"`
	OperatorName string      `json:"operator_name"`
}

var (
	// Table is the table name of order_review.
	Table       = "order_review"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "order_review ore"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
