package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

//综合征
var KlSyndrome = klSyndromeApi{}

type klSyndromeApi struct{}

// 获取
func (*klSyndromeApi) Tree(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	tree, err := service.KlSyndromeService.Tree(nil)
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}

func (*klSyndromeApi) Type(r *ghttp.Request) {
	tree := g.Array{}
	for _, v := range kl_syndrome.TypeTree {
		child := g.Array{}
		for _, v2 := range v.Children {
			child = append(child, g.Map{
				"value": v2.SubType,
				"label": v2.Name,
			})
		}

		tree = append(tree, g.Map{
			"value":    v.Type,
			"label":    v.Name,
			"children": child,
		})
	}
	response.Success(r, tree)
}

//添加医院
func (*klSyndromeApi) Add(r *ghttp.Request) {
	var req *kl_syndrome.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlSyndromeService.Add(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	var info *kl_syndrome.Entity
	kl_syndrome.M.Where(g.Map{"type": req.Type[0], "name": req.Name, "name_en": req.NameEn}).Scan(&info)
	response.SuccessMsg(r, "添加成功", g.Map{"id": info.Id})
}

func (*klSyndromeApi) Edit(r *ghttp.Request) {
	var req *kl_syndrome.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlSyndromeService.Edit(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "编辑成功")
}

// 同级排序
func (*klSyndromeApi) SetSort(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	type_str := r.GetQueryInt("type")
	sub_type_str := r.GetQueryInt("sub_type")
	compare_id := r.GetQueryInt("compare_id")
	t := r.GetQueryString("t")

	var list []*kl_feature.Entity
	if err := kl_syndrome.M.Where(g.Map{"type": type_str, "sub_type": sub_type_str}).Order("sort,id").Scan(&list); err != nil {
		response.ErrorDb(r, err)
	}

	id_arr := []int{}
	for _, v := range list {
		if v.Id != id {
			if v.Id == compare_id {
				if t == "after" {
					id_arr = append(id_arr, v.Id)
					id_arr = append(id_arr, id)
				} else if t == "before" {
					id_arr = append(id_arr, id)
					id_arr = append(id_arr, v.Id)
				}
			} else {
				id_arr = append(id_arr, v.Id)
			}
		}
	}

	for i, v := range id_arr {
		if _, err := kl_syndrome.M.Where("id", v).Data(g.Map{"sort": i}).Update(); err != nil {
			response.ErrorDb(r, err)
		}
	}

	response.Success(r)
}

//设置状态
func (*klSyndromeApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlSyndromeService.SetStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除医院
func (*klSyndromeApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlSyndromeService.Delete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

func (*klSyndromeApi) Details(r *ghttp.Request) {
	id := r.GetQueryInt("id")

	info, err := service.KlSyndromeService.Info(id)
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info != nil {
		var syndrome_feature_list []*kl_syndrome_feature.Entity
		err = kl_syndrome_feature.M.Where("syndrome_id", id).Order("id").Scan(&syndrome_feature_list)
		info.SyndromeFeatures = syndrome_feature_list
	}
	response.Success(r, info)
}
