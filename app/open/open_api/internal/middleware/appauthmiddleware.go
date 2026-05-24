package middleware

import (
	"net/http"

	"gorm.io/gorm"
)

type AppAuthMiddleware struct {
	DB *gorm.DB
}

func NewAppAuthMiddleware(db *gorm.DB) *AppAuthMiddleware {
	return &AppAuthMiddleware{
		DB: db,
	}
}

// Handle 应用级认证中间件（验证 appId + appSecret）
func (m *AppAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: 从请求体或 header 中获取 appId 和 appSecret，验证应用合法性
		// 这里先简单放行，后续根据实际需求实现

		next(w, r)
	}
}
