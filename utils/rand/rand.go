package utils

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// 生成指定长度的随机字符串
func GenerateUUId() string {
	return uuid.New().String()
}

// 随机生成器中的字符集合
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}
