package auth_leader_bind

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table auth_leader_key.
type Entity struct {
	Id                 int         `orm:"id,primary,size:4,table_comment:'主任秘钥绑定记录'" json:"id"`
	AuthLeaderId       int         `orm:"auth_leader_id,size:4,comment:'主任秘钥表主键'" json:"auth_leader_id"`
	SerialNumber       string      `orm:"serial_number,size:20,comment:'设备系列号'" json:"serial_number"`
	HospitalId         int         `orm:"hospital_id,size:4,comment:'医院ID'" json:"hospital_id"`
	ClientAuthorNumber string      `orm:"client_author_number,size:20,comment:'客户端授权码'" json:"client_author_number"`
	ServerAuthorNumber string      `orm:"server_author_number,size:20,comment:'服务端授权码'" json:"server_author_number"`
	CreateAt           *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt           *gtime.Time `orm:"update_at" json:"update_at"`
	Status             int         `orm:"status,size:2,not null,default:1" json:"status"`
	HospitalName       string      `json:"hospital_name"`
}

var (
	// Table is the table name of auth_leader_key.
	Table       = "auth_leader_bind"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_leader_bind alb"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
