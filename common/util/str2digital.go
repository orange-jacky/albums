package util

import "strconv"

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func StrToInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func StrToFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func Float64ToStr(s float64) string {
	return strconv.FormatFloat(s, 'f', -1, 64)
}
