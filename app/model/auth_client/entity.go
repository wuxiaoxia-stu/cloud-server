package auth_client

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id               int         `orm:"id,primary,size:4,table_comment:'授权客户端基本信息记录表'" json:"id"`
	Role             int         `orm:"role,size:2,not null" json:"role"`                                                                        // 客户端角色 （1-windows 客户端，2-web 端）
	Uuid             string      `orm:"uuid,size:40,not null" json:"uuid" p:"uuid" v:"required#参数错误,UUID必填"`                                     // 主板uuid
	UkeySerialNumber string      `orm:"ukey_serial_number,size:20" json:"ukey_serial_number" p:"ukey_serial_number" v:"required#参数错误,ukey系列号必填"` // 授权ukey系列号
	ClientVersion    string      `orm:"client_version,size:50" json:"client_version" p:"client_version" v:"required#参数错误,客户端版本号必填"`              // 客户端版本
	AiVersion        string      `orm:"ai_version,size:50" json:"ai_version"`                                                                    // ai版本
	CpuName          string      `orm:"cpu_name,size:100" json:"cpu_name"`                                                                       // cpu_name
	CpuID            string      `orm:"cpu_id,size:100" json:"cpu_id"`                                                                           // cpu_processor_id
	BaseboardID      string      `orm:"baseboard_id,size:100" json:"baseboard_id"`                                                               // baseboard_id
	GpuName          string      `orm:"gpu_name,size:200" json:"gpu_name"`                                                                       // gpu
	DiskID           string      `orm:"disk_id,size:300" json:"disk_id"`                                                                         // 硬盘id,几个拼起来
	CreateAt         *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt         *gtime.Time `orm:"update_at" json:"update_at"`
	Status           int         `orm:"status,size:2,not null,default:1" json:"status"`
	Signature        string      `json:"signature" p:"signature" v:"required#参数错误,签名值必填"`
	SerialNumber     string      `json:"serial_number"`
	HospitalId       int         `json:"hospital_id" p:"hospital_id"`
}

var (
	Table       = "auth_client"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_client ac"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type SignData struct {
	Uuid             string `json:"uuid"`
	SerialNumber     string `json:"serial_number"`
	UkeySerailNumber string `json:"ukey_serial_number"`
	CpuName          string `json:"cpu_name"`
	CpuID            string `json:"cpu_id""`
	BaseboardID      string `json:"baseboard_id"`
}
