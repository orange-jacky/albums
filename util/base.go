package util

import (
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