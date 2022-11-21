package middleware

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gookit/color"
)

// CORS 跨域处理中间件
func CORS(r *ghttp.Request) {
	color.Greenln("Router: ", r.GetUrl())
	if g.Cfg().GetBool("server.Debug") {
		color.Print("<fg=FF0066>提交数据:</>")
		color.Greenln(string(r.GetBody()))
	}
	r.Response.CORSDefault()
	r.Middleware.Next()
}
