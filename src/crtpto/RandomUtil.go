package crtpto

import (
	"math/rand"
	"time"
	"strconv"
)

// 利用当前时间的UNIX时间戳作为seed
// 返回值：小于等于n的一个随机Int
func RandGeneratorInt(maxValue int) int {
	rand.Seed(time.Now().UnixNano())
	return  rand.Intn(maxValue)
}

// 返回值：小于等于N的一个String
func RandGeneratorString(maxValue int) string {
	rand.Seed(time.Now().UnixNano())
	return  strconv.Itoa(rand.Intn(maxValue))
}