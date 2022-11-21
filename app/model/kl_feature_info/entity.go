package kl_feature_info

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table kl_feature.
type Entity struct {
	Id         int         `orm:"id,primary,size:4,table_comment:'知识图谱-特征详情描述'" json:"id"`
	Define     string      `orm:"define,size:text,comment:'定义'" json:"define"`
	DefineEn   string      `orm:"define_en,size:text,comment:'定义'" json:"define_en"`
	Diagnose   string      `orm:"diagnose,size:text,comment:'超声诊断要点'" json:"diagnose"`
	DiagnoseEn string      `orm:"diagnose_en,size:text,comment:'超声诊断要点'" json:"diagnose_en"`
	Consult    string      `orm:"consult,size:text,comment:'预后咨询'" json:"consult"`
	ConsultEn  string      `orm:"consult_en,size:text,comment:'超声诊断要点'" json:"consult_en"`
	Other      string      `orm:"other,size:text,comment:'其他'" json:"other"`
	OtherEn    string      `orm:"other_en,size:text,comment:'其他'" json:"other_en"`
	CreateAt   *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status     int         `orm:"status,size:2,not null,default:1" json:"status"`
}

var (
	// Table is the table name of kl_feature.
	Table       = "kl_feature_info"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_feature_info kfi"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
