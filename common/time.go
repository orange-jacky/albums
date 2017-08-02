package util

import (
	"time"
)

func GetNano() int64 {
	return time.Now().UnixNano()
}
