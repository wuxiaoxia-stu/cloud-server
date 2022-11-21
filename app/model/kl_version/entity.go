package kl_version

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table kl_version.
type Entity struct {
	Id           int         `orm:"id,primary,size:4,table_comment:'知识图谱-版本'" json:"id"`
	Name         string      `orm:"name,size:50,comment:'版本名称'" json:"name"`
	DataPath     string      `orm:"data_path,size:250,comment:'打包数据路径'" json:"data_path"`
	CreateAt     *gtime.Time `orm:"create_at" json:"create_at"`
	OperatorId   int         `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status       int         `orm:"status,size:2,not null,default:1" json:"status"`
	DataFullPath string      `json:"data_full_path"`
	OperatorName string      `json:"operator_name"`
}

var (
	Table       = "kl_version"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_version kv"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//添加特征表单验证规则
type AddReq struct {
	Name string `p:"name" v:"required|length:1,50|check-kl-version-name#版本名称必填|版本名称不超过50个字符长度|版本名称重复"`
}

func init() {
	//自定义验证规则，检查版本号名称是否重复
	if err := gvalid.RegisterRule("check-kl-version-name", CheckerKlversionName); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查综合征中文名是否重复
func CheckerKlversionName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where(g.Map{"name": value}).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}
