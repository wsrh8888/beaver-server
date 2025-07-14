package middleware

import (
	"beaver/common/ajax"
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

		// 获取城市信息并添加到header
		cityName, err := ajax.GetCityByIP(ClientIP)
		if err != nil {
			cityName = "" // 默认返回"未知"
		}
		r.Header.Set("X-City-Name", cityName)

		fmt.Println("scheme:", scheme, "ClientIP:", ClientIP, "ClientHost:", originalHost, "CityName:", cityName)
		next(w, r.WithContext(ctx))
	}
}
