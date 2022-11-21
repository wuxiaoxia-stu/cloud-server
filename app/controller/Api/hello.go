package Api

import (
	"aiyun_cloud_srv/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"time"
)

var Hello = hello{}

type hello struct{}

// Index is a demonstration route handler for output "Hello World!".
func (*hello) Index(r *ghttp.Request) {
	g.Dump(time.Now().AddDate(0, 36, 0).Unix())
	response.Error(r, "12312")
}
