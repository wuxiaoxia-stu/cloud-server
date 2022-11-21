package sys_admin

import (
	"aiyun_cloud_srv/app/model/sys_role"
	"context"
	"encoding/json"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gvalid"
	"os"
)

// Entity is the golang structure for table sys_admin.
type Entity struct {
	Id           int         `orm:"id,primary,table_comment:'用户表'" json:"id"`
	RoleId       int         `orm:"role_id,size:2,not null,default:1,comment:'角色id'" json:"role_id"` //角色ID
	DepartmentId string      `orm:"department_id,size:30,comment:'部门id'" json:"department_id"`       //部门ID
	Username     string      `orm:"username,not null,size:50,comment:'用户名'" json:"username"`         //用户名
	Password     string      `orm:"password,size:50,comment:'密码'" json:"password"`                   //密码
	Salt         string      `orm:"salt,size:5,comment:'加密盐值'" json:"salt"`                          // 盐值
	Phone        string      `orm:"phone,size:20,comment:'手机号'" json:"phone"`                        //手机号
	Email        string      `orm:"email,size:30,comment:'邮箱'" json:"email"`                         //邮箱
	Avatar       string      `orm:"avatar,size:150,comment:'头像'" json:"avatar"`                      //头像
	CreateAt     *gtime.Time `orm:"create_at,comment:'创建时间'" json:"create_at"`
	UpdateAt     *gtime.Time `orm:"update_at,comment:'修改时间'" json:"update_at"`
	Status       int         `orm:"status,size:2,default:1,comment:'状态，1：启用，0：禁用'" json:"status" default:"1"`
	Token        string      `json:"accessToken"`
	Expires      int         `json:"expires"`
	RoleRule     string      `json:"role_rule"`
	RoleName     string      `json:"role_name"`
}

var (
	// Table is the table name of sys_admin.
	Table       = "sys_admin"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "sys_admin sa"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//登录数据校验
type LoginReq struct {
	Username string `json:"username" p:"username" v:"required#用户名必填"`
	Password string `json:"password" p:"password" v:"required#密码必填"`
	//Code     string `json:"code" p:"code" v:"required|integer|length:4,4#验证码必填|验证码错误|验证码错误"`
}

//设置角色参数
type SetRoleReq struct {
	Id     int `json:"id" p:"id" v:"required#参数错误"`
	RoleId int `json:"role_id" p:"role_id" v:"required|check-role#参数错误|参数错误"`
}

func init() {
	//自定义验证规则，检查type值是否合法
	if err := gvalid.RegisterRule("check-role", CheckerRole); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-username", CheckerUsername); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查role_id值是否合法
func CheckerRole(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *sys_role.Entity
	err := sys_role.M.Where("id", value).Where("status", 1).Scan(&info)
	if err != nil {
		return err
	}

	if info == nil {
		return gerror.New(message)
	}

	return nil
}

//自定义验证规则，检查username值是否合法
func CheckerUsername(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var info *Entity
	err := M.Where("username", value).Scan(&info)
	if err != nil {
		return err
	}

	if info != nil {
		return gerror.New("用户名已经存在")
	}

	return nil
}

//添加用户表单验证规则
type AddReq struct {
	RoleId       int      `p:"role_id" v:"required|check-role#角色必填|角色不存在"`
	Username     string   `p:"username" v:"required|length:5,20|check-username#用户名必填|用户名长度限定为5到20个字符|用户名已存在"`
	Password     string   `p:"password" v:"required|length:5,20#密码必填|用户名长度限定为5到20个字符"`
	DepartmentId []string `p:"department_id" v:"required#部门信息必填"`
	Email        string   `p:"email" v:"required|email#邮箱必填|邮箱格式异常"`
	Phone        string   `p:"phone" v:"phone#手机号格式异常"`
	Status       int      `p:"status" v:"in:1,0#状态异常" default:"1"`
}

//修改用户表单验证规则
type EditReq struct {
	Id           int      `p:"id" v:"required#ID必填"`
	RoleId       int      `p:"role_id" v:"required|check-role#角色必填|角色不存在"`
	Password     string   `p:"password" v:"length:5,20#用户名长度限定为5到20个字符"`
	DepartmentId []string `p:"department_id" v:"required#部门信息必填"`
	Email        string   `p:"email" v:"required|email#邮箱必填|邮箱格式异常"`
	Phone        string   `p:"phone" v:"phone#手机号格式异常"`
	Status       int      `p:"status" v:"in:1,0#状态异常" default:"1"`
}

type DepartmentTree struct {
	Value    string            `json:"value"`
	Label    string            `json:"label"`
	Children []*DepartmentTree `json:"children"`
}

func GetDepartmentTree() (tree []*DepartmentTree, err error) {
	file, err := os.Open("./data/department.json")
	if err != nil {
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		g.Log().Error(err)
		return
	}

	err = json.Unmarshal(buffer, &tree)
	if err != nil {
		return
	}

	return
}
