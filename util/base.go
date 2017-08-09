package util

import (
	"fmt"
	"runtime"
)

//目录分隔符
func DirSeg() string {
	var s string
	switch runtime.GOOS {
	case "darwin", "linux":
		s = "/"
	case "window":
		s = "\\"
	default:
		s = "/"
	}
	return s
}

func GetDir(user, album, time string) string {
	return fmt.Sprintf("filecache%s%s%s%s-%s", DirSeg(), user, DirSeg(), album, time)
}

//生成访问图片路由
func GetUrl() string {
	conf := Configure("")
	return fmt.Sprintf("%s:%s/%s", conf.Nginx.HostInter, conf.Nginx.Port, conf.Nginx.Router)
}
