package main

import (
	"aiyun_cloud_srv/app/model/kl_feature"
	"aiyun_cloud_srv/app/model/kl_feature_atlas"
	"aiyun_cloud_srv/app/model/kl_syndrome"
	"aiyun_cloud_srv/app/model/kl_syndrome_feature"
	"aiyun_cloud_srv/app/model/sys_admin"
	"aiyun_cloud_srv/app/service"
	"aiyun_cloud_srv/library/utils/curl"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_DeepEqual(t *testing.T) {
	admin, _ := service.SysAdminService.FindById(1)

	admin.CreateAt = nil
	admin.UpdateAt = nil
	admin2 := &sys_admin.Entity{
		Id:           1,
		RoleId:       1,
		DepartmentId: "PD,developmentEngineer",
		Username:     "admin",
		Password:     "315193380d2fedffa677b7fe236fadb5",
		Salt:         "sinb",
		Phone:        "18388178881",
		Email:        "778774780@qq.com",
		Avatar:       "",
		CreateAt:     nil,
		UpdateAt:     gtime.NewFromStrFormat("2022-06-29 17:59:47", "Y-m-d H:i:s"),
		Status:       1,
		Expires:      0,
	}
	//admin2.CreateAt = nil
	admin2.UpdateAt = nil
	g.Dump(admin2)
	t.Log(reflect.DeepEqual(admin, admin2))

	t.Log(admin.CreateAt == admin2.CreateAt)
}

type Msg struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

type Part struct {
	PartUuid   string `json:"part_uuid"`
	PartSerial string `json:"part_serial"`
	PartName   string `json:"part_name"`
	PartNameEn string `json:"part_name_en"`
}

type Position struct {
	PositionUuid   string `json:"position_uuid"`
	PositionSerial string `json:"position_serial"`
	PositionName   string `json:"position_name"`
	PositionNameEn string `json:"position_name_en"`
}

type Group struct {
	GroupUuid   string     `json:"group_uuid"`
	GroupName   string     `json:"group_name"`
	GroupNameEn string     `json:"group_name_en"`
	GroupSerial string     `json:"group_serial"`
	Features    []*Feature `json:"features"`
}

type Feature struct {
	FeatureUuid        string `json:"feature_uuid"`
	FeatureSerial      string `json:"feature_serial"`
	FeatureName        string `json:"feature_name"`
	FeatureNameEn      string `json:"feature_name_en"`
	PositionUuid       string `json:"position_uuid"`
	PartUuid           string `json:"part_uuid"`
	SyndromeFeatureOpt string `json:"syndrome_feature_opt"`
}

type FeatureAtlas struct {
	FeatureUuid      string `json:"feature_uuid"`
	LegendFilePath   string `json:"legend_file_path"`
	LegendOriginName string `json:"legend_origin_name"`
	LegendType       string `json:"legend_type"`
}

type FeatureDetails struct {
	FeatureConsult    string          `json:"feature_consult"`
	FeatureConsultEn  string          `json:"feature_consult_en"`
	FeatureDefine     string          `json:"feature_define"`
	FeatureDefineEn   string          `json:"feature_define_en"`
	FeatureDiagnose   string          `json:"feature_diagnose"`
	FeatureDiagnoseEn string          `json:"feature_diagnose_en"`
	FeatureName       string          `json:"feature_name"`
	FeatureNameEn     string          `json:"feature_name_en"`
	FeatureOther      string          `json:"feature_other"`
	FeatureOtherEn    string          `json:"feature_other_en"`
	FeatureSerial     string          `json:"feature_serial"`
	FeatureUuid       string          `json:"feature_uuid"`
	GroupUuid         string          `json:"group_uuid"`
	PartUuid          string          `json:"part_uuid"`
	PositionUuid      string          `json:"position_uuid"`
	FeatureLegends    []*FeatureAtlas `json:"feature_legends"`
}

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6IjE1ODIwMjcwNTEwIiwidXNlcm5hbWUiOiJheWphZG1pbiIsInVzZXJfbmFtZV9yZWFsIjoi6LaF57qn566h55CG5ZGYIiwiZGVwYXJ0bWVudCI6IiIsInV1aWQiOiI1MzlhYzJiOC04ZGYxLTRiM2YtOWM3ZS0zMDBhY2Q3Zjg1MjQiLCJ1aWQiOjEsInN0YXRlIjoyLCJleHAiOjE2NTgzOTgwNzksImlzcyI6Imdpbi1ibG9nIn0.TdZm-uhKvXSLJUSq_5N_Q7FchVlJrUzaGXYZuS_jc9E"

//知识图谱-特征数据初始化
func Test_KL_Feature(t *testing.T) {
	kl_feature.M.Delete("id > 0")
	kl_feature_atlas.M.Delete("id > 0")

	b2, _ := curl.Get("http://192.168.1.8:8005/api/v2.5/kl_feature/not_classify_features?search_name=&page=1&limit=1000", g.Map{},
		g.Map{"Authorization": token})

	var msg2 Msg
	json.Unmarshal(b2, &msg2)

	var feature_list []*Feature
	if err := gconv.Struct(msg2.Data["list"], &feature_list); err != nil {
		g.Log().Error(err.Error())
	}

	b, _ := curl.Get("http://192.168.1.8:8005/api/v2.5/kl_feature/part_list?page=1&limit=1000&sort_type=1&search_type=&search_name=", g.Map{},
		g.Map{"Authorization": token})

	var msg Msg
	json.Unmarshal(b, &msg)

	var part_list []*Part
	if err := gconv.Struct(msg.Data["list"], &part_list); err != nil {
		g.Log().Error(err.Error())
	}

	for _, v := range part_list {
		kl_feature.M.Data(g.Map{
			"uuid":    v.PartUuid,
			"pid":     0,
			"level":   1,
			"name":    v.PartName,
			"name_en": v.PartNameEn,
		}).Insert()

		var info *kl_feature.Entity
		kl_feature.M.Where(g.Map{
			"uuid":    v.PartUuid,
			"pid":     0,
			"level":   1,
			"name":    v.PartName,
			"name_en": v.PartNameEn,
		}).Order("id DESC").Scan(&info)

		for _, v2 := range feature_list {
			if v2.PartUuid == v.PartUuid && v2.PositionUuid == "" {
				feature_details(v2.FeatureUuid, token, v2, info.Id, 1)
			}
		}

		b, _ := curl.Get("http://192.168.1.8:8005/api/v2.5/kl_feature/part_positions", g.Map{"page": 1, "limit": 1000, "sort_type": 1, "part_uuid": v.PartUuid},
			g.Map{"Authorization": token})

		var msg Msg
		json.Unmarshal(b, &msg)

		var position_list []*Position
		if err := gconv.Struct(msg.Data["list"], &position_list); err != nil {
			g.Log().Error(err.Error())
		}

		for _, v2 := range position_list {
			kl_feature.M.Data(g.Map{
				"uuid":    v2.PositionUuid,
				"pid":     info.Id,
				"level":   2,
				"name":    v2.PositionName,
				"name_en": v2.PositionNameEn,
			}).Insert()

			var info2 *kl_feature.Entity
			kl_feature.M.Where(g.Map{
				"uuid":    v2.PositionUuid,
				"pid":     info.Id,
				"level":   2,
				"name":    v2.PositionName,
				"name_en": v2.PositionNameEn,
			}).Order("id DESC").Scan(&info2)

			b, _ := curl.Get("http://192.168.1.8:8005/api/v2.5/kl_feature/position_groups",
				g.Map{"page": 1, "limit": 1000, "sort_type": 1, "part_uuid": v.PartUuid, "position_uuid": v2.PositionUuid},
				g.Map{"Authorization": token})

			var msg Msg
			json.Unmarshal(b, &msg)

			var group_list []*Group
			if err := gconv.Struct(msg.Data["list"], &group_list); err != nil {
				g.Log().Error(err.Error())
			}

			for _, v3 := range group_list {
				kl_feature.M.Data(g.Map{
					"uuid":    v3.GroupUuid,
					"pid":     info2.Id,
					"level":   3,
					"name":    v3.GroupName,
					"name_en": v3.GroupNameEn,
				}).Insert()

				var info3 *kl_feature.Entity
				kl_feature.M.Where(g.Map{
					"uuid":    v3.GroupUuid,
					"pid":     info2.Id,
					"level":   3,
					"name":    v3.GroupName,
					"name_en": v3.GroupNameEn,
				}).Order("id DESC").Scan(&info3)
				//b, _ := curl.Get("http://192.168.1.8:8005/api/v2.5/kl_feature/group_features",
				//	g.Map{"page": 1, "limit": 1000, "sort_type": 1, "part_uuid": v.PartUuid, "position_uuid": v2.PositionUuid, "group_uuid": v3.GroupUuid},
				//	g.Map{"Authorization": token})
				//
				//var msg Msg
				//json.Unmarshal(b, &msg)
				//
				//var feature_list []*Feature
				//if err := gconv.Struct(msg.Data["list"], &feature_list); err != nil {
				//	g.Log().Error(err.Error())
				//}

				for _, v4 := range v3.Features {
					feature_details(v4.FeatureUuid, token, v4, info3.Id, 0)
				}
			}

			for _, v3 := range feature_list {
				if v3.PositionUuid == v2.PositionUuid {
					feature_details(v3.FeatureUuid, token, v3, info2.Id, 1)
				}
			}
		}
	}
}

func feature_details(feature_uuid, token string, f *Feature, pid int, invisible int) {
	b, _ := curl.Get("http://192.168.1.8:8005/api/v1/kl_feature/feature_detail",
		g.Map{"feature_uuid": feature_uuid},
		g.Map{"Authorization": token})

	var msg Msg
	json.Unmarshal(b, &msg)

	var feature_details *FeatureDetails
	if err := gconv.Struct(msg.Data, &feature_details); err != nil {
		g.Log().Error(err.Error())
	}

	var feature_info *kl_feature.Entity
	kl_feature.M.Where(g.Map{
		"uuid": feature_details.FeatureUuid,
	}).Scan(&feature_info)

	if feature_info != nil {
		return
	}

	kl_feature.M.Data(g.Map{
		"uuid":        feature_details.FeatureUuid,
		"pid":         pid,
		"level":       4,
		"invisible":   invisible,
		"name":        f.FeatureName,
		"name_en":     f.FeatureNameEn,
		"define":      feature_details.FeatureDefine,
		"define_en":   feature_details.FeatureDefineEn,
		"diagnose":    feature_details.FeatureDiagnose,
		"diagnose_en": feature_details.FeatureDiagnoseEn,
		"consult":     feature_details.FeatureConsult,
		"consult_en":  feature_details.FeatureConsultEn,
		"other":       feature_details.FeatureOther,
		"other_en":    feature_details.FeatureOtherEn,
		"operator_id": 1,
	}).Insert()

	var info *kl_feature.Entity
	kl_feature.M.Where("uuid", feature_details.FeatureUuid).Scan(&info)

	for _, v := range feature_details.FeatureLegends {
		t := 1
		if v.LegendType == "FL03" {
			t = 3
		} else if v.LegendType == "FL02" {
			t = 2
		}
		kl_feature_atlas.M.Data(g.Map{
			"feature_id":  info.Id,
			"type":        t,
			"url":         v.LegendFilePath,
			"name":        v.LegendOriginName,
			"operator_id": 1,
		}).Insert()
	}
}

type Syndrome struct {
	GeneLocation       string `json:"gene_location"`
	GeneLocationEn     string `json:"gene_location_en"`
	GeneticsDesc       string `json:"genetics_desc"`
	GeneticsDescEn     string `json:"genetics_desc_en"`
	SyndromeConsult    string `json:"syndrome_consult"`
	SyndromeConsultEn  string `json:"syndrome_consult_en"`
	SyndromeDiagnose   string `json:"syndrome_diagnose"`
	SyndromeDiagnoseEn string `json:"syndrome_diagnose_en"`
	SyndromeGroup      int    `json:"syndrome_group"`
	SyndromeName       string `json:"syndrome_name"`
	SyndromeNameEn     string `json:"syndrome_name_en"`
	SyndromeSerial     string `json:"syndrome_serial"`
	SyndromeType       int    `json:"syndrome_type"`
	SyndromeUuid       string `json:"syndrome_uuid"`
}

//知识图谱-综合征数据初始化
func Test_KL_Syndrome(t *testing.T) {
	kl_syndrome.M.Delete("id > 0")
	kl_syndrome_feature.M.Delete("id > 0")

	var syndrome_list []*Syndrome

	for i := 1; i < 10; i++ {
		b, _ := curl.Get(fmt.Sprintf("http://192.168.1.8:8005/api/v1/kl_syndrome/syndrome_list?page=%d&limit=20&sort_type=1&search_name=", i), g.Map{}, g.Map{"Authorization": token})
		var msg Msg
		json.Unmarshal(b, &msg)

		var syndrome_list_curr_page []*Syndrome
		if err := gconv.Struct(msg.Data["list"], &syndrome_list_curr_page); err != nil {
			g.Log().Error(err.Error())
		}

		syndrome_list = append(syndrome_list, syndrome_list_curr_page...)
	}

	not_found_count := 0
	for _, v := range syndrome_list {
		b2, _ := curl.Get(fmt.Sprintf("http://192.168.1.8:8005/api/v1/kl_syndrome/syndrome_detail?syndrome_uuid=%s", v.SyndromeUuid), g.Map{}, g.Map{"Authorization": token})
		var msg2 Msg
		json.Unmarshal(b2, &msg2)

		var syndrome *Syndrome
		if err := gconv.Struct(msg2.Data, &syndrome); err != nil {
			g.Log().Error(err.Error())
		}

		b, _ := curl.Get(fmt.Sprintf("http://192.168.1.8:8005/api/v1/kl_syndrome/morphology_list_origin?syndrome_uuid=%s", v.SyndromeUuid), g.Map{}, g.Map{"Authorization": token})
		var msg Msg
		json.Unmarshal(b, &msg)

		var feature []*Feature
		if err := gconv.Struct(msg.Data["list"], &feature); err != nil {
			g.Log().Error(err.Error())
		}

		feature_ids := []string{}
		syndrome_feature_list := []*kl_syndrome_feature.Entity{}
		for _, v2 := range feature {
			var feature_info *kl_feature.Entity
			kl_feature.M.Where("uuid", v2.FeatureUuid).Order("id").Scan(&feature_info)
			if feature_info == nil {
				not_found_count = not_found_count + 1
				g.Dump(fmt.Sprintf("未查询到特征信息,总数量:%d", not_found_count))
				g.Dump(v2.FeatureUuid)
			} else {
				feature_ids = append(feature_ids, strconv.Itoa(feature_info.Id))

				t := ""
				if v2.SyndromeFeatureOpt == "SY01" {
					t = "1"
				} else if v2.SyndromeFeatureOpt == "SY02" {
					t = "2"
				}
				syndrome_feature_list = append(syndrome_feature_list, &kl_syndrome_feature.Entity{
					Type:      t,
					FeatureId: feature_info.Id,
				})
			}
		}

		_, err := kl_syndrome.M.Data(g.Map{
			"uuid":             syndrome.SyndromeUuid,
			"type":             syndrome.SyndromeType,
			"sub_type":         syndrome.SyndromeGroup,
			"name":             syndrome.SyndromeName,
			"name_en":          syndrome.SyndromeNameEn,
			"gene_location":    syndrome.GeneLocation,
			"gene_location_en": syndrome.GeneLocationEn,
			"genetics_desc":    syndrome.GeneticsDesc,
			"genetics_desc_en": syndrome.GeneticsDescEn,
			"diagnose":         syndrome.SyndromeDiagnose,
			"diagnose_en":      syndrome.SyndromeDiagnoseEn,
			"consult":          syndrome.SyndromeConsult,
			"consult_en":       syndrome.SyndromeConsultEn,
			"feature_ids":      strings.Join(feature_ids, ","),
			"operator_id":      1,
		}).Insert()
		if err != nil {
			g.Log().Error(err)
		}

		var info *kl_syndrome.Entity
		err = kl_syndrome.M.Where(g.Map{"uuid": syndrome.SyndromeUuid}).Scan(&info)
		if err != nil {
			g.Log().Error(err)
		}
		g.Dump(info)

		syndrome_feature_list_ := g.Array{}
		for _, v2 := range syndrome_feature_list {
			feature_id_path, _ := service.KlFeatureService.FullPath(v2.FeatureId, "")
			feature_root_path_arr := strings.Split(feature_id_path, ",")
			syndrome_feature_list_ = append(syndrome_feature_list_, g.Map{
				"syndrome_id":     info.Id,
				"feature_id":      v2.FeatureId,
				"type":            v2.Type,
				"feature_root_id": feature_root_path_arr[0],
				"feature_id_path": feature_id_path,
			})
		}
		kl_syndrome_feature.M.Data(syndrome_feature_list_).Insert()

	}

}
