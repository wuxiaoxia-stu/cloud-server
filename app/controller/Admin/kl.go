package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var Kl = klApi{}

type klApi struct{}

type QuerySyndromeReq struct {
	Id []int `p:"id"`
}

type SyndromeRsp struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Invisible int    `json:"invisible"`
}

type SyndromeCount struct {
	SyndromeId int `json:"syndrome_id"`
	Count      int `json:"count"`
}

// 特征树
func (*klApi) FeatureTree(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	tree, err := service.KlFeatureService.Tree(g.Map{"level <=": 3, "status": 1, "invisible": 0})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}

// 查询遗传综合征
func (*klApi) QuerySyndrome(r *ghttp.Request) {
	var req *QuerySyndromeReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	var feature_list []*kl_feature.Entity
	if err := kl_feature.M.Fields("id").WhereIn("pid", req.Id).Order("pid,id").Scan(&feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	feature_ids := []int{}
	for _, v := range feature_list {
		feature_ids = append(feature_ids, v.Id)
	}

	//SELECT syndrome_id,count(*) FROM "lpm_kl_syndrome_feature" WHERE feature_id IN (267,646)  GROUP BY "syndrome_id" ORDER BY "count" DESC
	var syndrome_count_list []*SyndromeCount
	if err := kl_syndrome_feature.M.
		Fields("syndrome_id,count(*) as count").
		WhereIn("feature_id", feature_ids).
		Group("syndrome_id").Order("count DESC").
		Scan(&syndrome_count_list); err != nil {
		response.ErrorDb(r, err)
	}

	syndrome_ids := []int{}
	for _, v := range syndrome_count_list {
		syndrome_ids = append(syndrome_ids, v.SyndromeId)
	}

	var syndrome_feature_list []*kl_syndrome_feature.Entity
	if err := kl_syndrome_feature.M_alias.
		Fields("ksf.*,kf.name as feature_name").
		LeftJoin("kl_feature kf", "kf.id = ksf.feature_id").
		WhereIn("ksf.syndrome_id", syndrome_ids).
		Scan(&syndrome_feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	var syndrome_list []*kl_syndrome.Entity
	if err := kl_syndrome.M.
		Fields("id,name").
		WhereIn("id", syndrome_ids).
		Scan(&syndrome_list); err != nil {
		response.ErrorDb(r, err)
	}

	var list []*kl_syndrome.Entity
	for _, v := range syndrome_ids {
		for _, v2 := range syndrome_list {
			if v == v2.Id {
				list = append(list, v2)
			}
		}
	}

	for _, v := range list {
		for _, v2 := range syndrome_feature_list {
			if v2.SyndromeId == v.Id {
				v.Features = append(v.Features, &kl_feature.Entity{
					Id:      v2.FeatureId,
					Name:    v2.FeatureName,
					Type:    v2.Type,
					IsCheck: utils.InArray(v2.FeatureId, feature_ids),
				})
			}
		}
	}

	response.Success(r, list)
}

// 综合征列表
func (*klApi) SyndromeTree(r *ghttp.Request) {
	tree, err := service.KlSyndromeService.Tree(g.Map{"status": 1})
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}

// 综合征详情
func (*klApi) SyndromeDetails(r *ghttp.Request) {
	id := r.GetQueryString("id")

	var info *kl_syndrome.Entity
	if err := kl_syndrome.M.Where(g.Map{"id": id, "status": 1}).Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Success(r)
	}

	var syndrome_feature_list []*kl_syndrome_feature.Entity
	if err := kl_syndrome_feature.M_alias.
		Fields("ksf.*,kf.name as feature_name,kf2.name as feature_root_name").
		LeftJoin("kl_feature kf", "kf.id = ksf.feature_id").
		LeftJoin("kl_feature kf2", "kf2.id = ksf.feature_root_id").
		Where(g.Map{"ksf.syndrome_id": id}).
		Scan(&syndrome_feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	var root_feature_list []*kl_feature.Entity
	if err := kl_feature.M.Fields("id,name").Where(g.Map{"pid": 0, "status": 1}).Order("sort,id").Scan(&root_feature_list); err != nil {
		response.ErrorDb(r, err)
	}

	for _, v := range root_feature_list {
		for _, v2 := range syndrome_feature_list {
			if v.Id == v2.FeatureRootId {
				v.Children = append(v.Children, &kl_feature.Entity{
					Id:   v2.FeatureId,
					Name: v2.FeatureName,
					Type: v2.Type,
				})
			}
		}
	}

	info.Features = root_feature_list
	response.Success(r, info)
}

// 综合征详情
func (*klApi) FeatureDetails(r *ghttp.Request) {
	id := r.GetQueryString("id")

	var info *kl_feature.Entity
	if err := kl_feature.M.Where(g.Map{"id": id, "status": 1}).Scan(&info); err != nil {
		response.ErrorDb(r, err)
	}

	if info == nil {
		response.Success(r)
	}

	var atlas_list []*kl_feature_atlas.Entity
	if err := kl_feature_atlas.M.Where("feature_id", id).Order("id").Scan(&atlas_list); err != nil {
		response.ErrorDb(r, err)
	}

	info.Atlas = atlas_list
	response.Success(r, info)
}
