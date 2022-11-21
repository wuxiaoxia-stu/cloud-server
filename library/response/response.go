package response

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gmutex"
	"github.com/gookit/color"
)

const (
	SuccessCode int = 0
	ErrorCode   int = 1
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

var (
	response = new(Response)
	mu       = gmutex.New()
)

//成功返回JSON
func Success(r *ghttp.Request, data ...interface{}) {
	response.RJson(r, SuccessCode, "ok", data...)
}

//成功返回JSON
func SuccessMsg(r *ghttp.Request, msg string, data ...interface{}) {
	response.RJson(r, SuccessCode, msg, data...)
}

//失败返回JSON
func Error(r *ghttp.Request, msg string, data ...interface{}) {
	response.RJson(r, ErrorCode, msg, data...)
}

//自定义Code返回信息
func Json(r *ghttp.Request, code int, msg string, data ...interface{}) {
	response.RJson(r, code, msg, data...)
}

//系统错误信息
func ErrorSys(r *ghttp.Request, err error) {
	g.Log().Error(err)
	response.RJson(r, ErrorCode, "系统异常")
}

//数据库错误信息
func ErrorDb(r *ghttp.Request, err error) {
	g.Log().Error(err)
	response.RJson(r, ErrorCode, "数据库异常")
}

// 标准返回结果数据结构封装。
// 返回固定数据结构的JSON:
// code:  状态码(200:成功,302跳转，和http请求状态码一至);
// msg:  请求结果信息;
// data: 请求结果,根据不同接口返回结果的数据结构不同;
func (res *Response) RJson(r *ghttp.Request, code int, msg string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	response = &Response{
		Code: code,
		Msg:  msg,
		Data: responseData,
	}
	r.SetParam("apiReturnRes", response)
	r.Response.WriteJson(response)
	if g.Cfg().GetBool("server.Debug") {
		b, _ := json.Marshal(response)
		color.Print("<fg=FF0066>响应数据:</>")
		color.Greenln(string(b))
	}

	r.Exit()
}
