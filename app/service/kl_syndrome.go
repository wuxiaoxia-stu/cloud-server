package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"strconv"
	"strings"
)

//综合征
var KlSyndromeService = new(klSyndromeService)

type klSyndromeService struct{}

//获取列表数据
func (s *klSyndromeService) Info(syndrome_id int) (res *kl_syndrome.Entity, err error) {
	err = kl_syndrome.M.Where("id", syndrome_id).Limit(1).Scan(&res)
	return
}

//获取列表数据
func (s *klSyndromeService) List(where interface{}, order string) (res []*kl_syndrome.Entity, err error) {
	err = kl_syndrome.M.Where(where).Order(order).Scan(&res)
	return
}

//获取当前节点及其所有子节点数据
func (s *klSyndromeService) Tree(where interface{}) (tree []*kl_syndrome.Entity, err error) {
	var list []*kl_syndrome.Entity

	err = kl_syndrome.M.Order("sort,id").Where(where).Scan(&list)
	if err != nil {
		return
	}

	data, _ := json.Marshal(kl_syndrome.TypeTree)
	if err = json.Unmarshal(data, &tree); err != nil {
		return
	}

	for _, v := range tree {
		for _, v2 := range list {
			if v2.SubType == 0 && v.Type == v2.Type {
				v.Children = append(v.Children, v2)
			}
		}
	}

	for _, v := range tree {
		for _, v2 := range v.Children {
			for _, v3 := range list {
				if v3.SubType > 0 && v3.SubType == v2.SubType {
					v2.Children = append(v2.Children, v3)
				}
			}
		}
	}

	return
}

func (s *klSyndromeService) Add(req *kl_syndrome.AddReq, operator_id int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	sub_type := 0
	if len(req.Type) == 2 {
		sub_type = req.Type[1]
	}

	feature_ids := []string{}
	for _, v := range req.Features {
		for _, v2 := range v.Children {
			feature_ids = append(feature_ids, strconv.Itoa(v2.Id))
		}
	}

	data := g.Map{
		"type":             req.Type[0],
		"sub_type":         sub_type,
		"name":             req.Name,
		"name_en":          req.NameEn,
		"gene_location":    req.GeneLocation,
		"gene_location_en": req.GeneLocationEn,
		"genetics_desc":    req.GeneticsDesc,
		"genetics_desc_en": req.GeneticsDescEn,
		"diagnose":         req.Diagnose,
		"diagnose_en":      req.DiagnoseEn,
		"consult":          req.Consult,
		"consult_en":       req.ConsultEn,
		"operator_id":      operator_id,
	}

	if len(feature_ids) > 0 {
		data["feature_ids"] = "," + strings.Join(feature_ids, ",") + ","
	}

	if _, err := tx.Model(kl_syndrome.Table).Data(data).Insert(); err != nil {
		tx.Rollback()
		return err
	}

	var syndrome_info *kl_syndrome.Entity
	if err = tx.Model(kl_syndrome.Table).Where(data).Scan(&syndrome_info); err != nil {
		tx.Rollback()
		return err
	}

	syndrome_feature_data := g.Array{}
	for _, v := range req.Features {
		for _, v2 := range v.Children {
			feature_root_path, err := KlFeatureService.FullPath(v2.Id, "")
			if err != nil {
				tx.Rollback()
				return err
			}
			syndrome_feature_data = append(syndrome_feature_data, g.Map{
				"syndrome_id":       syndrome_info.Id,
				"feature_id":        v2.Id,
				"type":              v2.Type,
				"feature_root_id":   v.Id,
				"feature_root_path": feature_root_path,
			})
		}
	}

	if len(syndrome_feature_data) > 0 {
		if _, err = tx.Model(kl_syndrome_feature.Table).Data(syndrome_feature_data).Insert(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *klSyndromeService) Edit(req *kl_syndrome.EditReq, operator_id int) (err error) {
	tx, err := g.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	//sub_type := 0
	//if len(req.Type) == 2 {
	//	sub_type = req.Type[1]
	//}

	data := g.Map{
		//"type":             req.Type[0],
		//"sub_type":         sub_type,
		"name":             req.Name,
		"name_en":          req.NameEn,
		"gene_location":    req.GeneLocation,
		"gene_location_en": req.GeneLocationEn,
		"genetics_desc":    req.GeneticsDesc,
		"genetics_desc_en": req.GeneticsDescEn,
		"diagnose":         req.Diagnose,
		"diagnose_en":      req.DiagnoseEn,
		"consult":          req.Consult,
		"consult_en":       req.ConsultEn,
		"operator_id":      operator_id,
	}

	feature_ids := []string{}
	for _, v := range req.Features {
		for _, v2 := range v.Children {
			feature_ids = append(feature_ids, strconv.Itoa(v2.Id))
		}
	}
	if len(feature_ids) > 0 {
		data["feature_ids"] = "," + strings.Join(feature_ids, ",") + ","
	}

	if _, err = tx.Model(kl_syndrome.Table).Where("id", req.Id).Data(data).Update(); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Model(kl_syndrome_feature.Table).Where("syndrome_id", req.Id).Delete(); err != nil {
		tx.Rollback()
		return err
	}

	syndrome_feature_data := g.Array{}
	for _, v := range req.Features {
		for _, v2 := range v.Children {
			feature_root_path, err := KlFeatureService.FullPath(v2.Id, "")
			if err != nil {
				tx.Rollback()
				return err
			}
			syndrome_feature_data = append(syndrome_feature_data, g.Map{
				"syndrome_id":       req.Id,
				"feature_id":        v2.Id,
				"type":              v2.Type,
				"feature_root_id":   v.Id,
				"feature_root_path": feature_root_path,
			})
		}
	}

	if len(syndrome_feature_data) > 0 {
		if _, err = tx.Model(kl_syndrome_feature.Table).Data(syndrome_feature_data).Insert(); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

//设置数据状态
func (s *klSyndromeService) SetStatus(req *model.SetStatusParams) (err error) {
	_, err = kl_syndrome.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *klSyndromeService) Delete(req *model.Ids) (err error) {
	_, err = kl_syndrome.M.WhereIn("id", req.Ids).Delete()
	return
}
