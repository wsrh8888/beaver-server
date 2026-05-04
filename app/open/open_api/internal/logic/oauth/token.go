package oauth

import (
	"time"

	"github.com/google/uuid"
)

// GenerateAccessToken 生成 access_token
func GenerateAccessToken() string {
	return uuid.New().String()
}

// GenerateRefreshToken 生成 refresh_token
func GenerateRefreshToken() string {
	return uuid.New().String()
}

// GetTokenExpiry 获取 token 过期时间戳(毫秒)
func GetTokenExpiry(expiresIn int64) int64 {
	return time.Now().Add(time.Duration(expiresIn) * time.Second).UnixMilli()
}
