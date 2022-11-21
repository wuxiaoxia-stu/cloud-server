package kl_syndrome

import (
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

// Entity is the golang structure for table kl_syndrome.
type Entity struct {
	Id               int                           `orm:"id,primary,size:4,table_comment:'知识图谱-综合征'" json:"id"`
	Uuid             string                        `orm:"uuid,size:50,comment:'UUID'" json:"uuid"`
	Type             int                           `orm:"type,size:2,comment:'类型:1：遗传综合征，2：宫内感染，3：致畸剂'" json:"type"`
	SubType          int                           `orm:"sub_type,size:2,comment:'子类型'" json:"sub_type"`
	Name             string                        `orm:"name,size:100,comment:'特征名称'" json:"name"`
	NameEn           string                        `orm:"name_en,size:100,comment:'特征英文名称'" json:"name_en"`
	GeneLocation     string                        `orm:"gene_location,size:50,comment:'基因点位'" json:"gene_location"`
	GeneLocationEn   string                        `orm:"gene_location_en,size:50,comment:'基因点位'" json:"gene_location_en"`
	GeneticsDesc     string                        `orm:"genetics_desc,size:50,comment:'遗传类型'" json:"genetics_desc"`
	GeneticsDescEn   string                        `orm:"genetics_desc_en,size:50,comment:'遗传类型'" json:"genetics_desc_en"`
	FeatureIds       string                        `orm:"feature_ids,size:text,comment:'特征id'" json:"feature_ids"`
	Diagnose         string                        `orm:"diagnose,size:text,comment:'超声诊断要点'" json:"diagnose"`
	DiagnoseEn       string                        `orm:"diagnose_en,size:text,comment:'超声诊断要点'" json:"diagnose_en"`
	Consult          string                        `orm:"consult,size:text,comment:'预后咨询'" json:"consult"`
	ConsultEn        string                        `orm:"consult_en,size:text,comment:'预后咨询'" json:"consult_en"`
	Other            string                        `orm:"other,size:text,comment:'其他'" json:"other"`
	OtherEn          string                        `orm:"other_en,size:text,comment:'其他'" json:"other_en"`
	Sort             int                           `orm:"sort,size:4,,default:0,comment:'排序'" json:"sort"`
	CreateAt         *gtime.Time                   `orm:"create_at" json:"create_at"`
	UpdateAt         *gtime.Time                   `orm:"update_at" json:"update_at"`
	OperatorId       int                           `orm:"operator_id,size:4,comment:'操作人'" json:"operator_id"`
	Status           int                           `orm:"status,size:2,not null,default:1" json:"status"`
	Children         []*Entity                     `json:"children"`
	Features         []*kl_feature.Entity          `json:"features"`
	SyndromeFeatures []*kl_syndrome_feature.Entity `json:"syndrome_features"`
}

var (
	// Table is the table name of kl_syndrome.
	Table       = "kl_syndrome"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_syndrome ks"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)

var TypeTree = []*Entity{
	&Entity{Type: 1, Name: "遗传综合征", Status: 1, Children: []*Entity{
		&Entity{Type: 1, SubType: 1, Name: "生长发育受限为特征", Status: 1},
		&Entity{Type: 1, SubType: 2, Name: "生长过度为特征", Status: 1},
		&Entity{Type: 1, SubType: 3, Name: "颜面部异常为特征", Status: 1},
		&Entity{Type: 1, SubType: 4, Name: "大脑异常为特征", Status: 1},
		&Entity{Type: 1, SubType: 5, Name: "肢体异常为特征", Status: 1},
		&Entity{Type: 1, SubType: 6, Name: "骨骼发育不良为特征", Status: 1},
		&Entity{Type: 1, SubType: 7, Name: "颅缝早闭为特征", Status: 1},
		&Entity{Type: 1, SubType: 8, Name: "多发异常为特征", Status: 1},
		&Entity{Type: 1, SubType: 9, Name: "软组织异常为特征", Status: 1},
		&Entity{Type: 1, SubType: 10, Name: "序列征和联合症", Status: 1},
		&Entity{Type: 1, SubType: 13, Name: "染色体异常综合征", Status: 1},
	}},
	&Entity{Type: 2, Name: "宫内感染", Status: 1},
	&Entity{Type: 3, Name: "致畸剂", Status: 1},
}

//添加综合征表单验证规则
type AddReq struct {
	Type []int `p:"type" v:"required#综合征类型必填"`
	//SubType        int
	Name           string `p:"name" v:"required|length:1,100|check-syndrome-name#综合征中文名称必填|综合征中文名称不超过100个字符长度|综合征中文名不能重复"`
	NameEn         string `p:"name_en" v:"required|length:1,100|check-syndrome-name-en#综合征英文名称必填|综合征英文名称不超过100个字符长度|综合征英文名不能重复"`
	GeneLocation   string
	GeneLocationEn string
	GeneticsDesc   string
	GeneticsDescEn string
	Diagnose       string
	DiagnoseEn     string
	Consult        string
	ConsultEn      string
	Features       []*Features `p:"features"`
}

type Features struct {
	Id       int         `p:"id"`
	Name     string      `p:"name"`
	Type     string      `p:"type"`
	Children []*Features `p:"features"`
}

//编辑综合征表单验证规则
type EditReq struct {
	Id int `p:"id" v:"required#ID值必填"`
	//Type           []int  `p:"type" v:"required#综合征类型必填"`
	Name           string `p:"name" v:"required|length:1,100|check-syndrome-name#综合征中文名称必填|综合征中文名称不超过100个字符长度|综合征中文名不能重复"`
	NameEn         string `p:"name_en" v:"required|length:1,100|check-syndrome-name-en#综合征英文名称必填|综合征英文名称不超过100个字符长度|综合征英文名不能重复"`
	GeneLocation   string
	GeneLocationEn string
	GeneticsDesc   string
	GeneticsDescEn string
	Diagnose       string
	DiagnoseEn     string
	Consult        string
	ConsultEn      string
	Features       []*Features `p:"features"`
}

func init() {
	//自定义验证规则，检查type值是否合法
	if err := gvalid.RegisterRule("check-syndrome-name", CheckerSyndromeName); err != nil {
		panic(err)
	}

	if err := gvalid.RegisterRule("check-syndrome-name-en", CheckerSyndromeNameEn); err != nil {
		panic(err)
	}

}

//自定义验证规则，检查综合征中文名是否重复
func CheckerSyndromeName(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var req *EditReq
	if err := gconv.Struct(data, &req); err != nil {
		g.Log().Error(err.Error())
		return err
	}

	where := g.Map{"name": value}
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
func CheckerSyndromeNameEn(ctx context.Context, rule string, value interface{}, message string, data interface{}) error {
	var req *EditReq
	if err := gconv.Struct(data, &req); err != nil {
		g.Log().Error(err.Error())
		return err
	}

	where := g.Map{"name_en": value}
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
