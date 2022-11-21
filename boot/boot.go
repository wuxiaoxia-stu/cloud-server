package boot

import (
	"aiyun_cloud_srv/app/model"
	"aiyun_cloud_srv/library/utils"
	_ "aiyun_cloud_srv/packed"
	"github.com/gogf/gf/frame/g"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	//清理日志
	clearLog()

	//初始化数据库
	model.DbInit()
}

//清理一个月之前的日志
func clearLog() {
	g.Log().Info("开始清理日志")

	log_path := []string{}
	log_path = append(log_path, g.Cfg().GetString("database.logger.Path"))
	log_path = append(log_path, g.Cfg().GetString("server.LogPath"))
	log_path = append(log_path, g.Cfg().GetString("logger.Path"))

	// 日志存活时间
	log_expire := g.Cfg().GetInt("logger.Expire", 30)

	for _, v := range log_path {
		path_exist, _ := utils.PathExists(v)
		if strings.Contains(v, "log/") && path_exist {
			filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
				if info.ModTime().Before(time.Now().AddDate(0, 0, -1*log_expire)) {
					if err := os.Remove(path); err != nil {
						g.Log().Errorf("日志文件【%s】删除失败", path)
					} else {
						g.Log().Infof("日志文件【%s】已删除", path)
					}
				}
				return nil
			})
		}
	}
	g.Log().Info("日志清理完成！")
}
