package hospital

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table auth_leader_key.
type Entity struct {
	Id          int         `orm:"id,primary,size:4,table_comment:'单位管理表'" json:"id"`
	Name        string      `orm:"name,size:20,comment:'单位名称'" json:"name"`
	LicenseCode string      `orm:"license_code,size:20,comment:'营业执照'" json:"license_code"`
	RegionId    string      `orm:"region_id,size:20,comment:'所属地区'" json:"region_id"`
	CreateAt    *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt    *gtime.Time `orm:"update_at" json:"update_at"`
	OperatorId  int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status      int         `orm:"status,size:2,not null,default:1" json:"status"`
	RegionName  string      `json:"region_name"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "hospital"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "hospital h"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type PageReqParams struct {
	Name      string   `p:"name"`
	RegionId  []string `p:"region_id"`
	Status    int      `p:"status" default:"-1"`
	Page      int      `p:"page" default:"1"`
	PageSize  int      `p:"page_size" default:"10"`
	Order     string   `p:"order" default:"id"`
	Sort      string   `p:"sort" default:"DESC"`
	StartTime string   `p:"start_time"`
	EndTime   string   `p:"end_time"`
}

//添加用户表单验证规则
type AddReq struct {
	Name        string   `p:"name" v:"required|length:1,20|check-name#单位名称必填|单位名称不超过20个字符长度|单位名称重复"`
	LicenseCode string   `p:"license_code" v:"size:18#营业执照格式异常"`
	RegionId    []string `p:"region_id" v:"required#所属区域必填"`
}

type EditReq struct {
	Id          int      `p:"id" v:"required#参数错误"`
	Name        string   `p:"name" v:"required|length:1,20|check-name#单位名称必填|单位名称不超过20个字符长度|单位名称重复"`
	LicenseCode string   `p:"license_code" v:"size:18#营业执照格式异常"`
	RegionId    []string `p:"region_id" v:"required#所属区域必填"`
}

func init() {
	if err := gvalid.RegisterRule("check-name", CheckerName); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查name值是否合法
func CheckerName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var req *EditReq
	if err := gconv.Struct(data, &req); err != nil {
		return err
	}

	where := g.Map{"name": value}
	if req.Id > 0 {
		where["id <> ?"] = req.Id
	}

	var info *Entity
	err := M.Where(where).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}
