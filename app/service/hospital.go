package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/hospital"
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

var HospitalService = new(hospitalService)

type hospitalService struct{}

func (s *hospitalService) Info(where interface{}) (res *hospital.Entity, err error) {
	err = hospital.M.Where(where).Limit(1).Scan(&res)
	return
}

func (s *hospitalService) Add(req *hospital.AddReq, operator_id int) (result sql.Result, err error) {
	result, err = hospital.M.Data(g.Map{
		"name":         req.Name,
		"license_code": req.LicenseCode,
		"region_id":    req.RegionId[2],
		"used_times":   0,
		"operator_id":  operator_id,
		"status":       1,
	}).Insert()
	return
}

func (s *hospitalService) Edit(req *hospital.EditReq, operator_id int) (result sql.Result, err error) {
	result, err = hospital.M.Where("id", req.Id).Data(g.Map{
		"name":         req.Name,
		"license_code": req.LicenseCode,
		"region_id":    req.RegionId[2],
		"used_times":   0,
		"operator_id":  operator_id,
		"status":       1,
	}).Update()
	return
}

//设置数据状态
func (s *hospitalService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = hospital.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *hospitalService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = hospital.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取分页列表
func (s *hospitalService) Page(req *hospital.PageReqParams) (total int, list []*hospital.Entity, err error) {
	M := hospital.M

	if req.Name != "" {
		M = M.WhereLike("name", "%"+req.Name+"%")
	}

	if len(req.RegionId) == 3 {
		M = M.Where("region_id=?", req.RegionId[2])
	}

	if req.Status != -1 {
		M = M.Where("status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order(req.Order + " " + req.Sort)
	}

	total, err = M.Group("id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*hospital.Entity, len(data))
	err = data.Structs(&list)
	if err != nil {
		return
	}

	for _, v := range list {
		region_name, _ := RegionService.GetNameById(v.RegionId)
		v.RegionName = region_name
	}

	return
}
