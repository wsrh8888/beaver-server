package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"beaver/utils/jwts"

	"github.com/go-redis/redis"
)

// VerifyWsToken 验证 WS 连接鉴权：JWT 签名 + Redis 登录态
// 防止被踢出/注销的 token 仍能建立 WS 连接
func VerifyWsToken(token, accessSecret, userID string, rdb *redis.Client) error {
	if token == "" {
		return errors.New("token不能为空")
	}

	claims, err := jwts.ParseToken(token, accessSecret)
	if err != nil {
		return errors.New("无效的token")
	}

	if claims.UserID != userID {
		return errors.New("token与用户不匹配")
	}

	// 检查 Redis 登录态：遍历所有设备类型，找到匹配的 token 即合法
	for _, dt := range []string{"desktop", "mobile", "web", "unknown"} {
		key := fmt.Sprintf("login_%s_%s", userID, dt)
		val, err := rdb.Get(key).Result()
		if err != nil || val == "" {
			continue
		}
		var info map[string]interface{}
		if json.Unmarshal([]byte(val), &info) != nil {
			continue
		}
		if storedToken, _ := info["token"].(string); storedToken == token {
			return nil
		}
	}

	return errors.New("登录已过期，请重新登录")
}
