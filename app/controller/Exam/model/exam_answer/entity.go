package exam_answer

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type Entity struct {
	Id         int         `orm:"id,primary,size:4,table_comment:'考核答题'" json:"id"`
	GroupHash  string      `orm:"group_hash,size:50,not null,comment:'组hash'" json:"group_hash"`
	Type       int         `orm:"type,size:2,default:0,comment:'0：用户答题，1：ai答题'" json:"type"`
	UId        int         `orm:"uid,not null" json:"uid"`
	Username   string      `orm:"username,size:20,not null" json:"username"`
	Paper      string      `orm:"paper,size:10,comment:'考卷'" json:"paper"`
	Subject    string      `orm:"subject,size:50,comment:'题目'" json:"subject"`
	Num        int         `orm:"num,comment:'题目数量'" json:"num"`
	UseTime    int         `orm:"use_time,comment:'用时'" json:"use_time"`
	Index      int         `orm:"index,comment:'答题进度'" json:"index"`
	Question   string      `orm:"question,size:250,comment:'问题'" json:"question"`
	Answer     string      `orm:"answer,size:100,comment:'正确答案'" json:"answer"`
	UserAnswer string      `orm:"user_answer,size:100,comment:'用户答案'" json:"user_answer"`
	IsCorrect  string      `orm:"is_correct,size:2,comment:'是否答对，0：错误 1：答对'" json:"is_correct"`
	CreateAt   *gtime.Time `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time `orm:"update_at" json:"update_at"`
	Status     int         `orm:"status,size:2,not null,default:1" json:"status"`
}

var (
	Table       = "exam_answer"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "exam_answer ea"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

type AnswerReq struct {
	UserAnswer string `v:"required#请选择一项答案"`
	UseTime    int    `v:"required#参数错误，答题时间必填"`
}
