package kl_syndrome_feature

import (
	"github.com/gogf/gf/frame/g"
)

// Entity is the golang structure for table kl_syndrome_feature.
type Entity struct {
	Id              int    `orm:"id,primary,size:4,table_comment:'知识图谱-综合征-特征-关系表'" json:"id"`
	SyndromeId      int    `orm:"syndrome_id,comment:'综合征ID'" json:"syndrome_id"`
	FeatureId       int    `orm:"feature_id,comment:'特征ID'" json:"feature_id"`
	Type            string `orm:"type,size:5,comment:'类型 1：特征病变，2：其它可见病变'" json:"type"`
	FeatureRootId   int    `orm:"feature_root_id,comment:'特征顶级节点ID'" json:"feature_root_id"`
	FeatureIdPath   string `orm:"feature_id_path,size:30,comment:'特征节点路径'" json:"feature_id_path"`
	FeatureName     string `json:"feature_name"`
	FeatureRootName string `json:"feature_root_name"`
}

var (
	// Table is the table name of kl_syndrome_feature.
	Table       = "kl_syndrome_feature"
	M           = g.DB("default").Model(Table).Safe()
	Table_alias = "kl_syndrome_feature ksf"
	M_alias     = g.DB("default").Model(Table_alias).Safe()
)
