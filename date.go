package gcode

import (
	"strings"
	"time"
)

// 获取当前日期
func Date() string {
	return time.Now().Format("2006-01-02")
}

// 获取当前时间
func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05-07")
}

// 获取年月日
func DateYearMonthDay() (year, month, day string) {
	str := Date()
	tmp := strings.Split(str, "-")
	return tmp[0], tmp[1], tmp[2]
}
