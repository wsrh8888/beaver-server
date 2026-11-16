// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package middleware

import (
	"net/http"

	"beaver/app/platform/platform_api/internal/config"
	corejwt "beaver/middleware/jwt"
	"github.com/zeromicro/go-zero/rest"
)

// JwtS2sMiddleware 由 .api 的 `@server middleware: JwtS2s` 声明接线。
// 验票实现在仓根 middleware/jwt，audience 与主仓业务面一致（core:api），这里只负责装配。
type JwtS2sMiddleware struct {
	mw rest.Middleware
}

func NewJwtS2sMiddleware(c config.Config) *JwtS2sMiddleware {
	return &JwtS2sMiddleware{mw: corejwt.NewJwtS2s(c.S2S.Secret, corejwt.AudienceCore)}
}

func (m *JwtS2sMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.mw(next)
}
