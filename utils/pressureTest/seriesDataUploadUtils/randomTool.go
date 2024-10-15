/*
  Package seriesDataUploadUtils
  @Author: Ahsen17
  @Github: https://github.com/Ahsen17
  @Time: 2024/10/16 0:34
  @Description: ...
*/

package seriesDataUploadUtils

import (
	"math/rand"
	"time"
)

var (
	strSeed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numSeed = time.Now().UnixNano()

	randSeed = rand.NewSource(numSeed)
)

func init() {
	rand.Seed(numSeed)
}

// randStr 生成指定长度范围内的随机字符串
func randStr(length int) string {
	bytes := []byte(strSeed)
	var result []byte
	r := rand.New(randSeed)
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// randStr 生成指定长度范围内的随机字符串，指定上下限
func randStrDoubleLimits(minLen int, maxLen int) string {
	length := randInt(minLen, maxLen)
	return randStr(length)
}

// randInt 生成指定范围内的随机整数
func randInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

// randFloat 生成指定范围内的随机浮点数
func randFloat(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
