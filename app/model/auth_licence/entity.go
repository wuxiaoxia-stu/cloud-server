package auth_licence

import (
	"aiyun_cloud_srv/app/model/auth_client"
	"aiyun_cloud_srv/app/model/order_device"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id                     int                  `orm:"id,primary,size:4,table_comment:'授权记录表'" json:"id"`
	LicenceId              int                  `orm:"licence_id,size:4,comment:'服务端授权id'" json:"licence_id"`                     // 服务端授权id
	AuthClientId           int                  `orm:"auth_client_id,not null,comment:'auth_client表id'" json:"auth_client_id"`    // 关联客户端表id
	Role                   int                  `orm:"role,not null,comment:'授权角色，1：客户端，2：服务端'" json:"role"`                      // 授权角色，1：客户端，2：服务端
	AuthorNumber           string               `orm:"author_number,unique,size:10,not null,comment:'授权编号'" json:"author_number"` // 授权编号
	UkeyCode               string               `orm:"ukey_code,size:5,comment:'ukey代码'" json:"ukey_code"`                        // ukey唯一编号
	UkeySerialNumber       string               `orm:"ukey_serial_number,size:20,comment:'ukey系列号'" json:"ukey_serial_number"`    // ukey唯一编号
	HospitalId             int                  `orm:"hospital_id,comment:'被授权单位'" json:"hospital_id"`                            // 被授权单位
	Licence                string               `orm:"licence,size:text,comment:'客户端验签信息'" json:"licence"`                        // ukey唯一编号
	PublicKey              string               `orm:"public_key,size:text,comment:'公钥'" json:"public_key"`                       // 授权公钥
	PrivateKey             string               `orm:"private_key,size:text,comment:'私钥'" json:"private_key"`                     // 授权私钥
	DeviceSerialNumber     string               `orm:"device_serial_number,size:50,comment:'设备系列号'" json:"device_serial_number"`  // 客户端序列号: 地址编号+医院编号+设备序列号
	Uuid                   string               `orm:"uuid,size:40" json:"uuid"`                                                  // 主板uuid
	CreateAt               *gtime.Time          `orm:"create_at" json:"create_at"`
	UpdateAt               *gtime.Time          `orm:"update_at" json:"update_at"`
	PairAt                 *gtime.Time          `orm:"pair_at,comment:'配对时间'" json:"pair_at"`
	Status                 int                  `orm:"status,size:2,not null,default:1,comment:'状态，1：有效， 0：无效'" json:"status"`
	ClientInfo             *auth_client.Entity  `json:"client_info"`
	DeviceInfo             *order_device.Entity `json:"device_info"`
	ClientVersion          string               `json:"client_version"`
	AiVersion              string               `json:"ai_version"`
	HospitalName           string               `json:"hospital_name"`
	ServerAuthorNumber     string               `json:"server_author_number"`
	PairCount              int                  `json:"pair_count"`                // 服务端配对次数
	PairClientAuthorNumber []string             `json:"pair_client_author_number"` // 配对客户端编号
}

var (
	Table       = "auth_licence"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_licence al"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type PageReqParams struct {
	UkeyCode     string   `p:"ukey_code"`
	RegionIds    []string `p:"region_ids"`
	HospitalId   int      `p:"hospital_id"`
	AuthorNumber string   `p:"author_number"`
	SerialNumber string   `p:"serial_number"`
	Uuid         string   `p:"uuid"`
	Role         int      `p:"role"`
	Status       int      `p:"status" default:"-1"`
	Page         int      `p:"page" default:"1"`
	PageSize     int      `p:"page_size" default:"10"`
	Order        string   `p:"order" default:"id"`
	Sort         string   `p:"sort" default:"DESC"`
	StartTime    string   `p:"start_time"`
	EndTime      string   `p:"end_time"`
}
