package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("读取请求体失败: %v", err)
		}
		// 恢复请求体
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// 记录请求信息
		logx.Infof("请求路径: %s, 方法: %s, 请求体: %s", r.URL.Path, r.Method, string(body))

		next(w, r)
	}
}
