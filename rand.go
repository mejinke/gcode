package gcode

import (
	"math/rand"
	"time"
)

//获取随机数
func Rand(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	s := max - min
	if s <= 0 {
		s = 1
	}
	return min + rand.Intn(s)
}
