package util

import (
	"time"
)

func GetMills() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetNano() int64 {
	return time.Now().UnixNano()
}
