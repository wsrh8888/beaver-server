package middleware

import (
	"net/http"
	"strings"
	"time"

	models "beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// AuthMiddleware OAuth2 认证中间件
func AuthMiddleware(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 1. 提取 Token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"code":401,"msg":"缺少认证令牌"}`, http.StatusUnauthorized)
				return
			}

			// 2. 验证 Bearer 格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"code":401,"msg":"无效的认证格式"}`, http.StatusUnauthorized)
				return
			}

			accessToken := parts[1]

			// 3. 验证 Token 有效性
			appID, err := validateAccessToken(db, accessToken)
			if err != nil {
				logx.Errorf("Token 验证失败: %v", err)
				http.Error(w, `{"code":401,"msg":"令牌无效或已过期"}`, http.StatusUnauthorized)
				return
			}

			// 4. 检查应用状态
			if !isAppActive(db, appID) {
				http.Error(w, `{"code":403,"msg":"应用已被禁用"}`, http.StatusForbidden)
				return
			}

			// 5. 注入 app_id 到 header（供后续 handler 使用）
			r.Header.Set("X-App-ID", appID)

			next(w, r)
		}
	}
}

// validateAccessToken 验证访问令牌
func validateAccessToken(db *gorm.DB, token string) (string, error) {
	var accessToken models.OpenAccessToken

	err := db.Where("access_token = ? AND status = ?", token, 1).First(&accessToken).Error
	if err != nil {
		return "", err
	}

	// 检查是否过期
	if accessToken.ExpiresAt > 0 && accessToken.ExpiresAt < getCurrentTimestamp() {
		// 标记为过期
		db.Model(&accessToken).Update("status", 2)
		return "", ErrTokenExpired
	}

	return accessToken.AppID, nil
}

// isAppActive 检查应用是否激活
func isAppActive(db *gorm.DB, appID string) bool {
	var app models.OpenApp

	err := db.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error
	return err == nil
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// 错误定义
var (
	ErrTokenExpired = &AuthError{Code: 401, Msg: "令牌已过期"}
	ErrTokenInvalid = &AuthError{Code: 401, Msg: "令牌无效"}
)

type AuthError struct {
	Code int
	Msg  string
}

func (e *AuthError) Error() string {
	return e.Msg
}
