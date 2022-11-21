package order

import (
	"aiyun_cloud_srv/app/model/auth_ukey"
	"aiyun_cloud_srv/app/model/hospital"
	"aiyun_cloud_srv/app/model/order_device"
	"aiyun_cloud_srv/app/model/order_review"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table order.
type Entity struct {
	Id                 int                    `orm:"id,primary,table_comment:'订单管理'" json:"id"`
	OrderNumber        string                 `orm:"order_number,unique,size:15,not null,comment:'订单号'" json:"order_number"`
	PrevOrderNumber    string                 `orm:"prev_order_number,size:15,comment:'升级前订单号'" json:"prev_order_number"`
	ContractNumber     string                 `orm:"contract_number,size:15,not null,comment:'合同号'" json:"contract_number" v:"required|length:15,15|check-contract-number#合同号必填|合同号应为15位整数|合同号重复"`
	HospitalId         int                    `orm:"hospital_id,not null,comment:'医院'" json:"hospital_id" v:"required|check-hospital#请选择单位|单位不存在"`
	UkeyCode           string                 `orm:"ukey_code,size:5,comment:'U-Key代码'" json:"ukey_code"`
	SaleName           string                 `orm:"sale_name,size:20,comment:'销售姓名'" json:"sale_name" v:"required#销售名称必填"`
	SalePhone          string                 `orm:"sale_phone,size:20,comment:'销售电话'" json:"sale_phone" v:"required#销售名电话必填"`
	Principal          string                 `orm:"principal,size:20,comment:'单位负责人'" json:"principal"`
	PrincipalPhone     string                 `orm:"principal_phone,size:20,comment:'单位负责人电话'" json:"principal_phone"`
	Contact            string                 `orm:"contact,size:20,comment:'单位联系人'" json:"contact"`
	ContactPhone       string                 `orm:"contact_phone,size:20,comment:'单位联系人电话'" json:"contact_phone"`
	Receiver           string                 `orm:"receiver,size:20,comment:'收货人'" json:"receiver"`
	ReceiverPhone      string                 `orm:"receiver_phone,size:20,comment:'收货人电话'" json:"receiver_phone"`
	ReceiverRegion     string                 `orm:"receiver_region,size:10,comment:'收货地址'" json:"receiver_region"`
	ReceiverAddress    string                 `orm:"receiver_address,size:50,comment:'收货详细地址'" json:"receiver_address"`
	ExpressType        int                    `orm:"express_type,size:2,comment:'发货方式'" json:"express_type"`
	ExpressNo          string                 `orm:"express_no,size:20,comment:'快递单号'" json:"express_no"`
	DeployRemark       string                 `orm:"deploy_remark,size:500,comment:'部署备注'" json:"deploy_remark"`
	DeployOperatorId   int                    `orm:"deploy_operator_id,comment:'撤销人',default:0" json:"deploy_operator_id"`
	DeployAt           *gtime.Time            `orm:"deploy_at,comment:'ukey绑定时间'" json:"deploy_at"`
	Count              int                    `orm:"count,size:2,comment:'设备数量',default:0" json:"count"`
	Maintenance        int                    `orm:"maintenance,size:2,comment:'维保月数',default:0" json:"maintenance"`
	Probation          int                    `orm:"probation,size:2,comment:'试用天数',default:0" json:"probation"`
	Remark             string                 `orm:"remark,size:500,comment:'备注'" json:"remark"`
	CreateAt           *gtime.Time            `orm:"create_at" json:"create_at"`
	UpdateAt           *gtime.Time            `orm:"update_at" json:"update_at"`
	OperatorId         int                    `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status             int                    `orm:"status,size:2,not null,default:1" json:"status"`
	ReceiverRegionName string                 `json:"receiver_region_name"`
	HospitalName       string                 `json:"hospital_name"`
	Hospital           *hospital.Entity       `json:"hospital"`
	Devices            []*order_device.Entity `json:"devices"`
	Reviews            []*order_review.Entity `json:"reviews"`
}

var (
	// Table is the table name of order.
	Table       = "order"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "order o"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type PageReqParams struct {
	KeyWord        string   `p:"keyword"`
	OrderNumber    string   `p:"order_number"`
	ContractNumber string   `p:"contract_number"`
	SerialNumber   string   `p:"serial_number"`
	RegionIds      []string `p:"region_ids"`
	HospitalId     int      `p:"hospital_id"`
	Status         int      `p:"status" default:"-1"`
	SetMeal        int      `p:"set_meal" default:"0"`
	UseType        int      `p:"use_type" default:"0"`
	Time           []string `p:"time"`
	Page           int      `p:"page" default:"1"`
	PageSize       int      `p:"page_size" default:"10"`
	Order          string   `p:"order" default:"id"`
	Sort           string   `p:"sort" default:"DESC"`
	StartTime      string   `p:"start_time"`
	EndTime        string   `p:"end_time"`
}

//订单审核
type ReviewReq struct {
	OrderNumber int    `p:"order_number" v:"required#参数错误"`
	Status      int    `p:"status" v:"required#参数错误"`
	Remark      string `p:"remark" v:"length:0,500#备注信息限制在500字符内"`
}

//订单审核
type RevokeReq struct {
	OrderNumber int    `p:"order_number" v:"required#参数错误"`
	Remark      string `p:"remark" v:"length:0,500#撤销理由限制在500字符内"`
}

//订单审核
type DeployReq struct {
	OrderNumber int    `p:"order_number" v:"required#参数错误"`
	UkeyCode    string `p:"ukey_code" v:"required|check-ukey-code#参数错误,U-Key信息必填|U-Key不存在"` // 授权ukey代号
	ExpressType int    `p:"express_type" v:"in:0,1#参数错误"`                                   // 发货类型 0自提 1物流
	ExpressNo   string `p:"express_no"`                                                     // 快递单号
	Remark      string `p:"remark" v:"length:0,500#撤销理由限制在500字符内"`
}

//订单升级
type UpgradeReq struct {
	OrderNumber    string                 `p:"order_number" v:"required#参数错误"`
	SaleName       string                 `p:"sale_name" v:"required#销售名称必填"`
	SalePhone      string                 `p:"sale_phone" v:"required#销售电话必填"`
	ContractNumber string                 `v:"required|length:15,15|check-contract-number#合同号必填|合同号应为15位整数|合同号重复"`
	Devices        []*order_device.Entity `p:"devices"`
	Remark         string                 `p:"remark" v:"length:0,500#撤销理由限制在500字符内"` // 备注
}

func init() {
	//自定义验证规则，检查合同号是否合法
	if err := gvalid.RegisterRule("check-contract-number", CheckerContractNumber); err != nil {
		panic(err)
	}

	//自定义验证规则，检查单位id值是否合法
	if err := gvalid.RegisterRule("check-hospital", CheckerHospital); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-ukey-code", CheckerUkeyCode); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查合同号值否合法
func CheckerContractNumber(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("contract_number", value).Where("status NOT IN(-1,7,8)").Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查单位id值否合法
func CheckerHospital(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *hospital.Entity
	err := hospital.M.Where("id", value).Where("status", 1).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查合同号值否合法
func CheckerUkeyCode(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *auth_ukey.Entity
	err := auth_ukey.M.Where("code", value).Where("status", 1).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}
