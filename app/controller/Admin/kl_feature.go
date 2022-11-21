package Admin

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

var KlFeature = klFeatureApi{}

type klFeatureApi struct{}

// 获取
func (*klFeatureApi) Tree(r *ghttp.Request) {
	var req *model.PageReqParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	//level <= 3 OR (invisible = 1 AND level = 4)
	tree, err := service.KlFeatureService.Tree("")
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, tree)
}

//获取特征分组
func (*klFeatureApi) GroupList(r *ghttp.Request) {
	group_list, err := service.KlFeatureService.GroupList()
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, group_list)
}

// 获取子集数据
func (*klFeatureApi) ChildList(r *ghttp.Request) {
	pid := r.GetQueryInt("pid")

	list, err := service.KlFeatureService.List(g.Map{"pid": pid}, "sort,id")
	if err != nil {
		response.ErrorDb(r, err)
	}

	response.Success(r, list)
}

//添加医院
func (*klFeatureApi) Add(r *ghttp.Request) {
	var req *kl_feature.AddReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlFeatureService.Add(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	var info *kl_feature.Entity
	kl_feature.M.Where(g.Map{"pid": req.Pid, "name": req.Name, "name_en": req.NameEn}).Scan(&info)
	response.SuccessMsg(r, "添加成功", g.Map{"id": info.Id})
}

func (*klFeatureApi) Edit(r *ghttp.Request) {
	var req *kl_feature.EditReq

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlFeatureService.Edit(req, r.GetCtxVar("uid").Int()); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "编辑成功")
}

// 同级排序
func (*klFeatureApi) SetSort(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	pid := r.GetQueryInt("pid")
	compare_id := r.GetQueryInt("compare_id")
	t := r.GetQueryString("t")

	var list []*kl_feature.Entity
	if err := kl_feature.M.Where("pid", pid).Order("sort,id").Scan(&list); err != nil {
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
		if _, err := kl_feature.M.Where("id", v).Data(g.Map{"sort": i}).Update(); err != nil {
			response.ErrorDb(r, err)
		}
	}

	response.Success(r)
}

//设置医院状态
func (*klFeatureApi) SetStatus(r *ghttp.Request) {
	var req *model.SetStatusParams

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlFeatureService.SetStatus(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//删除医院
func (*klFeatureApi) Delete(r *ghttp.Request) {
	var req *model.Ids

	if err := r.Parse(&req); err != nil {
		response.Error(r, err.(gvalid.Error).FirstString())
	}

	if err := service.KlFeatureService.Delete(req); err != nil {
		response.ErrorDb(r, err)
	}

	response.SuccessMsg(r, "删除成功")
}

//特征详情
func (*klFeatureApi) Details(r *ghttp.Request) {
	id := r.GetQueryInt("id")

	info, err := service.KlFeatureService.Info(g.Map{"id": id})
	if err != nil {
		response.ErrorDb(r, err)
	}

	if info != nil {
		atlas, err := service.KlFeatureService.GetAtlas(g.Map{"feature_id": id})
		if err != nil {
			response.ErrorDb(r, err)
		}
		info.Atlas = atlas
	}

	response.Success(r, info)
}
