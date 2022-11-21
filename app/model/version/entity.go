package version

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

// Entity is the golang structure for table version
type Entity struct {
	Id              int         `orm:"id,primary,table_comment:'设备管理'" json:"id"`
	Channel         int         `orm:"channel,size:2,not null,default:0,comment:'渠道'" v:"required#渠道必填" json:"channel"`        //渠道
	VersionNumber   string      `orm:"version_number,size:20,not null,comment:'版本号'" v:"required#版本号必填" json:"version_number"` //版本号
	Info            string      `orm:"info,size:1000,comment:'更新信息'" v:"required#更新信息必填" json:"info"`                          //更新信息
	BugFix          string      `orm:"bug_fix,size:1000,comment:'Bug修复'" json:"bug_fix"`                                       //bug修复
	UpdateRange     string      `orm:"update_range,size:150,comment:'更新范围'" json:"update_range"`                               //更新范围
	PackageUrl      string      `orm:"package_url,size:150,comment:'更新包路径'" v:"required#必须上传更新包" json:"package_url"`           //更新包路径
	Remark          string      `orm:"remark,size:500,comment:'备注'" v:"length:0,500#备注信息在500字符内" json:"remark"`                //备注
	OperatorId      int         `orm:"operator_id,not null,comment:'操作人'" json:"operator_id"`
	UpdateAt        *gtime.Time `orm:"update_at" json:"update_at"`
	Status          int         `orm:"status,size:2,default:1,comment:'状态'" json:"status"`
	OperatorName    string      `json:"operator_name"`
	PackageFullPath string      `json:"package_full_path"`
}

var (
	// Table is the table name of version.
	Table       = "version"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "version v"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
