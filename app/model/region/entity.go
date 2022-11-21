package region

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table auth_leader_key.
type Entity struct {
	Id         int         `orm:"id,primary,table_comment:'区域管理表'" json:"id"`
	Pid        int         `orm:"pid,comment:'父级id'" json:"pid"`
	Name       string      `orm:"name,size:20,comment:'区域名称'" json:"name"`
	CreateAt   *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status     int         `orm:"status,size:2,not null,default:1" json:"status"`
	Children   []*Entity   `json:"children"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "region"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "region r"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
