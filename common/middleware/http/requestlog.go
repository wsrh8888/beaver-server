package httpMiddleware

import (
	"beaver/common/middleware/utils"
	"bytes"
	"io"
	"net/http"
	"time"
)

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.LogRequest(r.Method, r.URL.Path, string(body), nil, err, startTime)
			return
		}
		// 恢复请求体
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		next(w, r)

		// 记录请求信息
		utils.LogRequest(r.Method, r.URL.Path, string(body), nil, nil, startTime)
	}
}
