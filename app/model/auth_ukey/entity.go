package auth_ukey

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table auth_ukey.
type Entity struct {
	Id           int         `orm:"id,primary,size:4,table_comment:'授权Ukey记录'" json:"id"`
	SerialNumber string      `orm:"serial_number,size:20,comment:'硬件ID'" json:"serial_number"`
	Code         string      `orm:"code,size:50,comment:'UKey代码'" json:"code"`
	Type         int         `orm:"type,size:2,comment:'UKey类型'" json:"type"`
	PublicKey    string      `orm:"public_key,size:500,comment:'公钥'" json:"public_key"`
	AadminId     int         `orm:"admin_id,size:4,comment:'公钥'" json:"admin_id"`
	AuthTimes    int         `orm:"auth_times,comment:'授权次数'" json:"auth_times"`
	UsedTimes    int         `orm:"used_times,comment:'已使用授权次数'" json:"used_times"`
	CreateAt     *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId   int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status       int         `orm:"status,size:2,not null,default:1" json:"status"`
	Username     string      `json:"username"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "auth_ukey"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_ukey au"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type PageReqParams struct {
	SerialNumber string `p:"serial_number"`
	Code         string `p:"code"`
	Type         int    `p:"type"`
	Username     string `p:"username"`
	Status       int    `p:"status" default:"-1"`
	Page         int    `p:"page" default:"1"`
	PageSize     int    `p:"page_size" default:"10"`
	Order        string `p:"order" default:"id"`
	Sort         string `p:"sort" default:"DESC"`
	StartTime    string `p:"start_time"`
	EndTime      string `p:"end_time"`
}

type AddReq struct {
	SerialNumber string `p:"serial_number" v:"required|size:16|check-SerialNumber#U-Key异常|U-Key异常|U-Key已被绑定"`
	Code         string `p:"code" v:"required|length:5,5|check-code#U-Key标识必填|U-Key标识错误|U-Key标识已经存在"`
	Type         int    `p:"type" v:"required|in:1,2,3#授权类型必填|授权类型错误"`
	Publickey    string `p:"public_key" v:"required#授权类型必填"`
	AadminId     int    `p:"admin_id" v:"required#绑定用户必填"`
	AuthTimes    int    `p:"auth_times" v:"required|integer|between:1,99#绑定用户必填|授权次数为正整数|授权次数最大值为99的整数"`
}

func init() {
	//自定义验证规则，检查code值是否合法
	if err := gvalid.RegisterRule("check-SerialNumber", CheckerSerialNumber); err != nil {
		panic(err)
	}

	//自定义验证规则，检查code值是否合法
	if err := gvalid.RegisterRule("check-code", CheckerCode); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查CheckerSerialNumber值是否存在
func CheckerSerialNumber(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var entity []*Entity
	err := M.Where("serial_number", value).Scan(&entity)
	if err != nil {
		return err
	}

	if entity != nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查code值是否存在
func CheckerCode(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var entity []*Entity
	err := M.Where("code", value).Scan(&entity)
	if err != nil {
		return err
	}

	if entity != nil {
		return gerror.New(message)
	}

	return nil
}

//检查ukey状态是否正常
type QueryUkeyStatus struct {
	Uuid         string `p:"uuid" v:"required#参数错误：uuid不存在"`
	SerialNumber string `p:"serial_number" v:"required#参数错误：ukey系列号不存在"`
	Signature    string `p:"signature" v:"required#参数错误：签名信息必填"`
}
