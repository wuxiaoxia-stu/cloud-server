package auth_licence_num

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id       int         `orm:"id,primary,size:4,table_comment:'授权编号生成'" json:"id"`
	LMonth   string      `orm:"l_month,not null,comment:'授权id'" json:"l_month"`
	Number   int         `orm:"number,size:4,not null,comment:'自增号'" json:"number"`
	UpdateAt *gtime.Time `orm:"update_at,comment:'更新时间'" json:"update_at"`
}

var (
	Table       = "auth_licence_num"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_licence_num aln"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
