package sys_role

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table sys_admin.
type Entity struct {
	Id       int         `orm:"id,primary,table_comment:'角色管理表'" json:"id"`
	Name     string      `orm:"name,size:50,not null,comment:'角色名称'" json:"name"`
	Rule     string      `orm:"rule,size:text,comment:'权限规则'" json:"rule"`
	CreateAt *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt *gtime.Time `orm:"update_at" json:"update_at"`
	Status   int         `orm:"status,size:2,not null,default:1" json:"status"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "sys_role"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "sys_role sr"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type AddReq struct {
	Name string   `v:"required|length:1,50|check-role-name#参数错误:角色名称必填|角色名称50字符长度内|角色重复"`
	Rule []string `v:"required#参数错误：权限规则必填"`
}

type EditReq struct {
	Id   int      `v:"required#参数错误:ID必填"`
	Rule []string `v:"required#参数错误：权限规则必填"`
}

func init() {
	//自定义验证规则，检查合同号是否合法
	if err := gvalid.RegisterRule("check-role-name", CheckerRoleName); err != nil {
		panic(err)
	}
}

//自定义验证规则，检查合同号值否合法
func CheckerRoleName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("name", value).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New(message)
	}

	return nil
}
