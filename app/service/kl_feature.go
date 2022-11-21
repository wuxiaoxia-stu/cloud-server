package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"strconv"
)

var KlFeatureService = new(klFeatureService)

type klFeatureService struct{}

//获取列表数据
func (s *klFeatureService) Info(where interface{}) (res *kl_feature.Entity, err error) {
	err = kl_feature.M.Where(where).Scan(&res)
	return
}

//获取列表数据
func (s *klFeatureService) GetAtlas(where interface{}) (res []*kl_feature_atlas.Entity, err error) {
	err = kl_feature_atlas.M.Where(where).Scan(&res)
	return
}

//获取列表数据
func (s *klFeatureService) List(where interface{}, order string) (res []*kl_feature.Entity, err error) {
	err = kl_feature.M.Where(where).Order(order).Scan(&res)
	return
}

//获取当前节点及其所有子节点数据
func (s *klFeatureService) Tree(where interface{}) (tree []*kl_feature.Entity, err error) {
	var list []*kl_feature.Entity

	//where := ""
	//if is_invisible {
	//	where = "level <= 2 OR (invisible = 1 AND level = 4)"
	//} else {
	//	where = "level <= 3 OR (invisible = 0 AND level = 4)"
	//}

	err = kl_feature.M.Order("level,sort,id").Where(where).Scan(&list)
	if err != nil {
		return
	}

	tree = s.ListToTree(list, 0, tree)
	return
}

func (s *klFeatureService) GroupList() (group_list []*kl_feature.GroupList, err error) {
	tree, err := s.Tree("level <= 3 OR (invisible = 1 AND level = 4)")
	if err != nil {
		return
	}

	var list []*kl_feature.Entity
	err = kl_feature.M.Fields("id,pid,name,name_en,invisible").Order("level,sort,id").Where("invisible = 0").Scan(&list)
	if err != nil {
		return
	}

	for _, v := range tree {
		var children = []*kl_feature.GroupList{}

		if v.Invisible == 1 {
			children = append(children, &kl_feature.GroupList{
				Id:   v.Id,
				Name: v.Name,
			})
		}
		for _, v2 := range v.Children {
			if v2.Invisible == 1 {
				children = append(children, &kl_feature.GroupList{
					Id:   v2.Id,
					Name: v2.Name,
				})
			}

			for _, v3 := range v2.Children {
				for _, v4 := range list {
					if v3.Id == v4.Pid {
						children = append(children, &kl_feature.GroupList{
							Id:   v4.Id,
							Name: v4.Name,
						})
					}
				}

				if v3.Invisible == 1 {
					children = append(children, &kl_feature.GroupList{
						Id:   v3.Id,
						Name: v3.Name,
					})
				}
			}
		}

		group_list = append(group_list, &kl_feature.GroupList{
			Id:       v.Id,
			Name:     v.Name,
			Check:    make([]int, 0),
			Children: children,
		})
	}

	return
}

//获取当前节点及其所有子节点数据
func (s *klFeatureService) ListToTree(source []*kl_feature.Entity, id int, result []*kl_feature.Entity) []*kl_feature.Entity {
	if source == nil {
		return result
	}
	if len(source) == 0 {
		return result
	}

	var otherData []*kl_feature.Entity //多余的数据，用于下一次递归
	for _, v := range source {
		if v.Pid == id {
			result = append(result, v)
		} else {
			otherData = append(otherData, v)
		}

	}

	if result != nil {
		for i := 0; i < len(result); i++ {
			result[i].Children = s.ListToTree(otherData, result[i].Id, result[i].Children)
		}
	}

	return result
}

func (s *klFeatureService) Add(req *kl_feature.AddReq, operator_id int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	level := 1
	if req.Pid > 0 {
		var info *kl_feature.Entity
		err = kl_feature.M.Where("id", req.Pid).Scan(&info)
		if err != nil {
			tx.Rollback()
			return
		}

		if info == nil {
			tx.Rollback()
			return fmt.Errorf("上级不存在")
		}

		level = info.Level + 1
	}

	if _, err = tx.Model(kl_feature.Table).Data(g.Map{
		"pid":         req.Pid,
		"level":       level,
		"name":        req.Name,
		"name_en":     req.NameEn,
		"invisible":   req.Invisible,
		"define":      req.Define,
		"define_en":   req.DefineEn,
		"diagnose":    req.Diagnose,
		"diagnose_en": req.DefineEn,
		"consult":     req.Consult,
		"consult_en":  req.ConsultEn,
		"other":       req.Other,
		"other_en":    req.OtherEn,
		"operator_id": operator_id,
	}).Insert(); err != nil {
		tx.Rollback()
		return
	}

	if len(req.ImgList) > 0 {
		var info *kl_feature.Entity
		if err = tx.Model(kl_feature.Table).Where(g.Map{
			"pid":         req.Pid,
			"level":       level,
			"name":        req.Name,
			"name_en":     req.NameEn,
			"operator_id": operator_id,
		}).Order("id DESC").Scan(&info); err != nil {
			tx.Rollback()
			return
		}

		feature_atlas_list := g.Array{}
		for _, v := range req.ImgList {
			feature_atlas_list = append(feature_atlas_list, g.Map{
				"feature_id": info.Id,
				"type":       v.Type,
				"url":        v.Path,
				"name":       v.Name,
				"ext":        v.Ext,
				"size":       v.Size,
			})
		}

		if _, err = tx.Model(kl_feature_atlas.Table).Insert(feature_atlas_list); err != nil {
			tx.Rollback()
			return
		}
	}

	return tx.Commit()
}

func (s *klFeatureService) Edit(req *kl_feature.EditReq, operator_id int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if _, err = tx.Model(kl_feature.Table).Where("id", req.Id).
		Data(g.Map{
			"name":        req.Name,
			"name_en":     req.NameEn,
			"define":      req.Define,
			"define_en":   req.DefineEn,
			"diagnose":    req.Diagnose,
			"diagnose_en": req.DiagnoseEn,
			"consult":     req.Consult,
			"consult_en":  req.ConsultEn,
			"other":       req.Other,
			"other_en":    req.OtherEn,
			"operator_id": operator_id,
		}).Update(); err != nil {
		tx.Rollback()
		return
	}

	if len(req.ImgList) > 0 {
		//删除旧的图片
		if _, err = tx.Model(kl_feature_atlas.Table).Delete("feature_id", req.Id); err != nil {
			tx.Rollback()
			return
		}

		//添加新图片
		feature_atlas_list := g.Array{}
		for _, v := range req.ImgList {
			feature_atlas_list = append(feature_atlas_list, g.Map{
				"feature_id": req.Id,
				"type":       v.Type,
				"url":        v.Path,
				"name":       v.Name,
				"ext":        v.Ext,
				"size":       v.Size,
			})
		}

		if _, err = tx.Model(kl_feature_atlas.Table).Insert(feature_atlas_list); err != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit()
}

//设置数据状态
func (s *klFeatureService) SetStatus(req *model.SetStatusParams) (err error) {
	_, err = kl_feature.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *klFeatureService) Delete(req *model.Ids) (err error) {
	_, err = kl_feature.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取id路径地址
func (s *klFeatureService) FullPath(id int, path string) (full_path string, err error) {
	full_path = path
	var info *kl_feature.Entity
	err = kl_feature.M.Where("id", id).Scan(&info)
	if err != nil {
		return path, err
	}

	if info != nil {
		if info.Pid > 0 {
			if full_path != "" {
				full_path = strconv.Itoa(info.Id) + "," + full_path
			} else {
				full_path = strconv.Itoa(info.Id)
			}
			return s.FullPath(info.Pid, full_path)
		} else {
			if full_path != "" {
				full_path = strconv.Itoa(info.Id) + "," + full_path
			} else {
				full_path = strconv.Itoa(info.Id)
			}
		}
	}

	return
}
