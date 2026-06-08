package utils

import (
	"time"

	"beaver/utils/logger"
	"beaver/utils/logger/model"
)

var accessLog = logger.New("http")

// LogRequest 记录 HTTP/gRPC 请求日志
func LogRequest(method, path string, req, resp interface{}, err error, startTime time.Time) {
	duration := time.Since(startTime)
	data := map[string]interface{}{
		"method":   method,
		"path":     path,
		"duration": duration.String(),
	}
	if req != nil {
		data["req"] = req
	}
	if resp != nil {
		data["resp"] = resp
	}

	if err != nil {
		data["err"] = err.Error()
		accessLog.Error(model.LogMsg{
			Text: "请求失败",
			Data: data,
		})
		return
	}

	accessLog.Info(model.LogMsg{
		Text: "请求成功",
		Data: data,
	})
}
