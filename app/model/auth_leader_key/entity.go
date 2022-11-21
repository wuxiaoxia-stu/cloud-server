package auth_leader_key

import (
	"aiyun_cloud_srv/app/model/hospital"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table auth_leader_key.
type Entity struct {
	Id                 int         `orm:"id,primary,size:4,table_comment:'授权主任秘钥记录'" json:"id"`
	SerialNumber       string      `orm:"serial_number,size:20,comment:'硬件ID'" json:"serial_number"`
	Code               string      `orm:"code,size:50,comment:'主任秘钥代码'" json:"code"`
	HospitalId         int         `orm:"hospital_id,size:4,comment:'医院ID'" json:"hospital_id"`
	AuthCount          int         `orm:"auth_count,size:2,comment:'授权数量',default:0" json:"auth_count"`
	Manager            string      `orm:"manager,size:20,comment:'管理员用户名称'" json:"manager"`
	Public             string      `orm:"public,size:20,comment:'普通用户名称'" json:"public"`
	CreateAt           *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt           *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId         int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status             int         `orm:"status,size:2,not null,default:1" json:"status"`
	HospitalName       string      `json:"hospital_name"`
	ServerAuthorNumber string      `json:"server_author_number"`
}

var (
	// Table is the table name of auth_leader_key.
	Table       = "auth_leader_key"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "auth_leader_key alk"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type AddReq struct {
	SerialNumber string `p:"serial_number" v:"required|size:16|check-leader-SerialNumber#主任秘钥异常|主任秘钥异常|主任秘钥已被注册"`
	Code         string `p:"code" v:"required|length:5,5|check-leader-code#主任秘钥标识必填|主任秘钥标识错误|主任秘钥标识已经存在"`
	HospitalId   string `p:"hospital_id" v:"required|check-HospitalId#单位必填|单位不存在或被禁用"`
	Manager      string `p:"manager" v:"required|length:1,20#管理员名称必填|管理员名称不超过20字符长度"`
	Public       string `p:"public" v:"required|length:1,20#普通用户名称必填|普通用户名称不超过20字符长度"`
}

func init() {
	//自定义验证规则，检查code值是否合法
	if err := gvalid.RegisterRule("check-leader-SerialNumber", CheckerSerialNumber); err != nil {
		panic(err)
	}

	//自定义验证规则，检查code值是否合法
	if err := gvalid.RegisterRule("check-leader-code", CheckerCode); err != nil {
		panic(err)
	}

	//自定义验证规则，检查HospitalId值是否合法
	if err := gvalid.RegisterRule("check-HospitalId", CheckeHospitalId); err != nil {
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

//自定义验证规则，检查code值是否存在
func CheckeHospitalId(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var entity []*Entity
	err := hospital.M.Where("id", value).Where("status", 1).Scan(&entity)
	if err != nil {
		return err
	}

	if entity == nil {
		return gerror.New(message)
	}

	return nil
}

//绑定主任秘钥提交数据
type BindLeaderReq struct {
	//LeaderKeyCode      string `json:"leader_key_code" p:"leader_key_code" v:"required#参数错误：主任秘钥代码必填"`           // 主任密钥代号
	LeaderSerialNumber string `json:"leader_serial_number" p:"leader_serial_number" v:"required#参数错误,主任秘钥系列号必填"` // 主任密钥唯一编号
	AuthorNumber       string `json:"author_number" p:"author_number" v:"required#参数错误：客户端授权码必填"`                // 客户端授权id
	ServerAuthorNumber string `json:"server_author_number" p:"server_author_number" v:"required#参数错误：服务端授权码必填"`  // 服务端授权id
	//Signature          string `json:"signature" p:"signature" v:"required#参数错误：数据签名必填"`                         // 数据签名
}

//绑定主任秘钥提交数据
type UnbindLeaderReq struct {
	LeaderSerialNumber string `json:"leader_serial_number" p:"leader_serial_number" v:"required#参数错误,主任秘钥系列号必填"` // 主任密钥唯一编号
	ServerAuthorNumber string `json:"server_author_number" p:"server_author_number" v:"required#参数错误：服务端授权码必填"`  // 服务端授权id
}
