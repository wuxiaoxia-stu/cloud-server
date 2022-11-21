package Admin

import (
	"aiyun_cloud_srv/library/response"
	"aiyun_cloud_srv/library/utils"
	"fmt"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
)

var File = fileApi{}

type fileApi struct{}

func (*fileApi) Upload(r *ghttp.Request) {
	tempFile := r.GetUploadFile("file")
	scene := r.GetFormString("scene")
	if scene == "" {
		response.Error(r, "参数错误")
	}

	if tempFile == nil {
		response.Error(r, "未获取到文件")
	}

	//保存文件
	date := gtime.Date()
	path := fmt.Sprintf("%s/%s/", scene, date)
	ext := utils.Ext(tempFile.Filename)
	fileName := tempFile.Filename
	tempFile.Filename = gmd5.MustEncrypt(grand.Letters(10)) + "." + ext

	_, err := tempFile.Save("public/" + path)
	if err != nil {
		response.Error(r, "文件上传失败")
	}

	base_url := g.Cfg().GetString("server.Domain")
	response.Success(r, g.Map{
		"base_url": base_url,
		"path":     path + tempFile.Filename,
		"name":     fileName,
		"size":     tempFile.Size,
		"ext":      ext,
	})
}
