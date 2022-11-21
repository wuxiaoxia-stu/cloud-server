package service

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/app/model/auth_ukey"
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

var AuthUkeyService = new(authUkeyService)

type authUkeyService struct{}

//查询数据
func (s *authUkeyService) Find(where g.Map) (list []*auth_ukey.Entity, err error) {
	err = auth_ukey.M.Where(where).Scan(&list)
	return
}

func (s *authUkeyService) FindOne(where g.Map) (data *auth_ukey.Entity, err error) {
	err = auth_ukey.M.Where(where).Limit(1).Scan(&data)
	return
}

func (s *authUkeyService) Add(req *auth_ukey.AddReq, operator_id int) (result sql.Result, err error) {
	result, err = auth_ukey.M.Data(g.Map{
		"serial_number": req.SerialNumber,
		"code":          req.Code,
		"type":          req.Type,
		"public_key":    req.Publickey,
		"admin_id":      req.AadminId,
		"auth_times":    req.AuthTimes,
		"used_times":    0,
		"operator_id":   operator_id,
		"status":        1,
	}).Insert()
	return
}

//设置数据状态
func (s *authUkeyService) SetStatus(req *model.SetStatusParams) (result sql.Result, err error) {
	result, err = auth_ukey.M.Where("id", req.Id).Data(req).Update()
	return
}

//批量删除
func (s *authUkeyService) Delete(req *model.Ids) (result sql.Result, err error) {
	result, err = auth_ukey.M.WhereIn("id", req.Ids).Delete()
	return
}

//获取分页列表
func (s *authUkeyService) Page(req *auth_ukey.PageReqParams) (total int, list []*auth_ukey.Entity, err error) {
	M := auth_ukey.M_alias

	if req.Username != "" {
		M = M.WhereLike("sa.username", "%"+req.Username+"%")
	}

	if req.SerialNumber != "" {
		M = M.WhereLike("au.serial_number", "%"+req.SerialNumber+"%")
	}

	if req.Code != "" {
		M = M.WhereLike("au.code", "%"+req.Code+"%")
	}

	if req.Type > 0 {
		M = M.Where("au.type=?", req.Type)
	}

	if req.Status != -1 {
		M = M.Where("au.status=?", req.Status)
	}

	if req.StartTime != "" {
		M = M.WhereGTE("au.create_at", req.StartTime)
	}
	if req.EndTime != "" {
		M = M.WhereLTE("au.create_at", req.EndTime)
	}
	if req.Order != "" && req.Sort != "" {
		M = M.Order("au." + req.Order + " " + req.Sort)
	}

	total, err = M.LeftJoin("sys_admin sa", "sa.id = au.admin_id").Group("au.id").Count()
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取总行数失败")
		return
	}

	M = M.Page(req.Page, req.PageSize)
	data, err := M.Fields("au.*,sa.username").LeftJoin("sys_admin sa", "sa.id = au.admin_id").All()

	if err != nil {
		g.Log().Error(err)
		err = gerror.New("获取数据失败")
		return
	}
	list = make([]*auth_ukey.Entity, len(data))
	err = data.Structs(&list)
	return
}
