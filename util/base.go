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

func GetDir(user, album string) string {
	return fmt.Sprintf("filecache%s%s%s%s", DirSeg(), user, DirSeg(), album)
}
