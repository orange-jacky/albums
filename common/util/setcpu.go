package util

import (
	"runtime"
)

func SetupCPU() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)
}
