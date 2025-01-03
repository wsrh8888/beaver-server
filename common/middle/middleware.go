package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LogMiddleware 自定义的中间件
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ClientIP := httpx.GetRemoteAddr(r)
		ctx := context.WithValue(r.Context(), "ClientIP", ClientIP)

		// 获取请求中的原始域名信息
		originalHost := r.Host
		ctx = context.WithValue(ctx, "ClientHost", originalHost)

		// 判断请求协议是否为 HTTPS 并记录 "http" 或 "https"
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		ctx = context.WithValue(ctx, "Scheme", scheme)
		fmt.Println("scheme:", scheme, "ClientIP:", ClientIP, "ClientHost:", originalHost)
		next(w, r.WithContext(ctx))
	}
}
