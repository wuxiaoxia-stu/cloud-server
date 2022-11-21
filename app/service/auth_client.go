package service

import (
	"aiyun_cloud_srv/app/model/auth_client"
	"github.com/gogf/gf/frame/g"
)

var AuthClientService = new(authClientService)

type authClientService struct{}

//查询数据
func (s *authClientService) FindByUUID(uuid string, role int) (res *auth_client.Entity, err error) {
	err = auth_client.M.Where(g.Map{"uuid": uuid, "role": role, "status": 1}).Order("id DESC").Limit(1).Scan(&res)
	return
}
