package kl_feature

import (
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table kl_feature.
type Entity struct {
	Id         int                        `orm:"id,primary,size:4,table_comment:'知识图谱-特征'" json:"id"`
	Pid        int                        `orm:"pid,size:4,comment:'上级ID'" json:"pid"`
	Uuid       string                     `orm:"uuid,size:50,comment:'UUID'" json:"uuid"`
	Level      int                        `orm:"level,size:2,comment:'层级'" json:"level"`
	Name       string                     `orm:"name,size:100,comment:'特征名称'" json:"name"`
	NameEn     string                     `orm:"name_en,size:100,comment:'特征英文名称'" json:"name_en"`
	Invisible  int                        `orm:"invisible,size:2,default:0,comment:'超声不可见病理特征'" json:"invisible"`
	Define     string                     `orm:"define,size:text,comment:'定义'" json:"define"`
	DefineEn   string                     `orm:"define_en,size:text,comment:'定义'" json:"define_en"`
	Diagnose   string                     `orm:"diagnose,size:text,comment:'超声诊断要点'" json:"diagnose"`
	DiagnoseEn string                     `orm:"diagnose_en,size:text,comment:'超声诊断要点'" json:"diagnose_en"`
	Consult    string                     `orm:"consult,size:text,comment:'预后咨询'" json:"consult"`
	ConsultEn  string                     `orm:"consult_en,size:text,comment:'超声诊断要点'" json:"consult_en"`
	Other      string                     `orm:"other,size:text,comment:'其他'" json:"other"`
	OtherEn    string                     `orm:"other_en,size:text,comment:'其他'" json:"other_en"`
	Sort       int                        `orm:"sort,size:4,,default:0,comment:'排序'" json:"sort"`
	CreateAt   *gtime.Time                `orm:"create_at" json:"create_at"`
	UpdateAt   *gtime.Time                `orm:"update_at" json:"update_at"`
	OperatorId int                        `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status     int                        `orm:"status,size:2,not null,default:1" json:"status"`
	Children   []*Entity                  `json:"children"`
	Atlas      []*kl_feature_atlas.Entity `json:"atlas"`
	Type       string                     `json:"type"`     // 1:特征病变 2:其它可见病变
	IsCheck    bool                       `json:"is_check"` // 选中
}

var (
	// Table is the table name of kl_feature.
	Table       = "kl_feature"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_feature kf"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

//添加特征表单验证规则
type AddReq struct {
	Name       string `p:"name" v:"required|length:1,50|check-feature-name#特征中文名称必填|特征中文名称不超过50个字符长度|特征中文名称重复"`
	NameEn     string `p:"name_en" v:"required|length:1,50|check-feature-name-en#特征英文名称必填|特征英文名称不超过50个字符长度|特征英文名称重复"`
	Pid        int    `p:"pid" v:"required#PID必填"`
	Level      int
	Invisible  int
	Define     string
	DefineEn   string
	Diagnose   string
	DiagnoseEn string
	Consult    string
	ConsultEn  string
	Other      string
	OtherEn    string
	ImgList    []*ImgList
}

type ImgList struct {
	Ext  string
	Name string
	Path string
	Size int
	Type int
}

//编辑特征表单验证规则
type EditReq struct {
	Id         int `p:"id" v:"required#ID必填"`
	Pid        int
	Name       string `p:"name" v:"required|length:1,50|check-feature-name#特征中文名称必填|特征中文名称不超过50个字符长度|特征中文名称重复"`
	NameEn     string `p:"name_en" v:"required|length:1,50|check-feature-name-en#特征英文名称必填|特征英文名称不超过50个字符长度|特征英文名称重复"`
	Define     string
	DefineEn   string
	Diagnose   string
	DiagnoseEn string
	Consult    string
	ConsultEn  string
	Other      string
	OtherEn    string
	ImgList    []*ImgList
}

type GroupList struct {
	Id       int          `json:"value"`
	Name     string       `json:"label"`
	Check    []int        `json:"check"`
	Type     int          `json:"type"`
	Children []*GroupList `json:"children"`
}

func init() {
	//自定义验证规则，检查type值是否合法
	if err := gvalid.RegisterRule("check-feature-name", CheckerFeatureName); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-feature-name-en", CheckerFeatureNameEn); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查综合征中文名是否重复
func CheckerFeatureName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var req *EditReq
	if err := gconv.Struct(data, &req); err != nil {
		g.Log().Error(err.Error())
		return err
	}

	where := g.Map{"name": value, "pid": req.Pid}
	if req.Id > 0 {
		where["id !="] = req.Id
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

//自定义验证规则，检查综合征英文名是否重复
func CheckerFeatureNameEn(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var req *EditReq
	if err := gconv.Struct(data, &req); err != nil {
		g.Log().Error(err.Error())
		return err
	}

	where := g.Map{"name": value, "pid": req.Pid}
	if req.Id > 0 {
		where["id !="] = req.Id
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
