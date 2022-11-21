package exam_user

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id       int         `orm:"id,primary,size:4,table_comment:'考核用户'" json:"id"`
	Username string      `orm:"username,size:20,not null" json:"username"`
	Token    string      `orm:"token,size:text" json:"token"`
	CreateAt *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt *gtime.Time `orm:"update_at" json:"update_at"`
	Status   int         `orm:"status,size:2,not null,default:1" json:"status"`
}

var (
	Table       = "exam_user"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "exam_user eu"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type LoginReq struct {
	Username  string `v:"required|length:1,20#用户名必填|用户名不超过20个字符长度"`
	Password  string `v:"required#密码必填"`
	PaperType string `v:"required#考卷必填"`
	Subject   string
	Num       int `default:0`
	AppType   int
	Restart   bool
}
