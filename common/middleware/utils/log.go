package utils

import (
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogRequest 记录请求日志
func LogRequest(method, path string, req, resp interface{}, err error, startTime time.Time) {
	duration := time.Since(startTime)

	if err != nil {
		logx.Errorf("请求失败: %s %s, 耗时: %v, 请求: %v, 错误: %v",
			method, path, duration, req, err)
		return
	}

	logx.Infof("请求成功: %s %s, 耗时: %v, 请求: %v, 响应: %v",
		method, path, duration, req, resp)
}
