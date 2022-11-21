package router

import (
	"aiyun_cloud_srv/app/controller/Admin"
	"aiyun_cloud_srv/app/controller/Api"
	"aiyun_cloud_srv/app/controller/Exam"
	"aiyun_cloud_srv/middleware"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()

	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		group.ALL("/public", Admin.Public)

		group.Middleware(middleware.Auth)
		group.ALL("/file", Admin.File)
		group.ALL("/index", Admin.Index)
		group.ALL("/order", Admin.Order)
		group.ALL("/order/manager_review", Admin.Order.Review) //权限问题 接口单独提出配置路由
		group.ALL("/licence", Admin.Licence)
		group.ALL("/server/list", Admin.Licence.List) //权限问题 接口单独提出配置路由
		group.ALL("/client/list", Admin.Licence.List) //权限问题 接口单独提出配置路由
		group.ALL("/sys-admin", Admin.SysAdmin)
		group.ALL("/sys-role", Admin.SysRole)
		group.ALL("/ukey", Admin.AuthUKey)
		group.ALL("/leader-key", Admin.LeaderKey)
		group.ALL("/hospital", Admin.Hospital)
		group.ALL("/version", Admin.Version)
		group.ALL("/kl-feature", Admin.KlFeature)
		group.ALL("/kl-syndrome", Admin.KlSyndrome)
		group.ALL("/kl-version", Admin.KlVersion)

	})

	s.Group("/api/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)

		//测试使用
		group.ALL("kl", Admin.Kl)

		group.ALL("public", Api.Public)

		group.ALL("exam", Exam.Index).Middleware(Exam.Auth)

		group.ALL("/hello", Api.Hello)
		group.ALL("/hospital", Api.Hospital)
		group.ALL("/authorize", Api.Authorize)
		group.ALL("/licence", Api.Licence)
		group.ALL("/pair", Api.Pair)
		group.ALL("/leader-key", Api.LeaderKey)
	})

}
