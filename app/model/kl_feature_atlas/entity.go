package kl_feature_atlas

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table kl_feature_info.
type Entity struct {
	Id         int         `orm:"id,primary,size:4,table_comment:'知识图谱-特征图集'" json:"id"`
	FeatureId  int         `orm:"feature_id,size:4,comment:'特征信息表ID'" json:"feature_id"`
	Type       int         `orm:"type,size:2,comment:'类型：0：典型超声图，1：病例解刨图，2：动态视频'" json:"type"`
	Url        string      `orm:"url,size:250,comment:'资源地址'" json:"url"`
	Name       string      `orm:"name,size:150,comment:'名称'" json:"name"`
	Ext        string      `orm:"ext,size:10,comment:'文件扩展名'" json:"ext"`
	Size       int         `orm:"size,comment:'文件大小'" json:"ext"`
	Sort       int         `orm:"sort,size:4,,default:0,comment:'排序'" json:"sort"`
	CreateAt   *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status     int         `orm:"status,size:2,not null,default:1" json:"status"`
}

var (
	// Table is the table name of kl_feature.
	Table       = "kl_feature_atlas"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_feature_atlas kfa"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
