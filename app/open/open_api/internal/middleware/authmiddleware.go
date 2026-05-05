package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"beaver/app/open/open_models"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthMiddleware struct {
	DB *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		DB: db,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 从 Header 中获取 Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":   401,
				"msg":    "缺少认证信息",
				"result": nil,
			})
			return
		}

		// 2. 解析 Bearer Token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":   401,
				"msg":    "无效的认证格式",
				"result": nil,
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":   401,
				"msg":    "Token 不能为空",
				"result": nil,
			})
			return
		}

		// 3. 查询 Token 是否有效
		var accessToken open_models.OpenAccessToken
		if err := m.DB.Where("token = ?", token).First(&accessToken).Error; err != nil {
			logx.Errorf("Token 查询失败: %v", err)
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":   401,
				"msg":    "无效的 Token",
				"result": nil,
			})
			return
		}

		// 4. 检查 Token 是否过期
		if time.Now().Unix() > accessToken.ExpiresAt {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":   401,
				"msg":    "Token 已过期",
				"result": nil,
			})
			return
		}

		// 5. 将 AppID 和 UserID 注入到 context 中
		ctx := context.WithValue(r.Context(), "appID", accessToken.AppID)
		if accessToken.UserID != "" {
			ctx = context.WithValue(ctx, "userID", accessToken.UserID)
		}

		// 6. 继续处理请求
		next(w, r.WithContext(ctx))
	}
}
