package ua

import (
	"beaver/utils/device"
	"context"
	"net/http"
)

// Context Key 定义
const (
	KeyDeviceID    = "ua-device-id"    // 设备唯一标识（从请求头获取）
	KeyDeviceType  = "ua-device-type"  // 精准设备类型：windows/macos/linux/ios/android
	KeyDeviceGroup = "ua-device-group" // 设备族群：desktop/mobile（用于互踢）
)

// Middleware UA 识别中间件
// 从 User-Agent 中提取设备信息并注入 Context
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		deviceID := r.Header.Get("DeviceID") // 从请求头获取设备ID
		preciseType := device.GetDeviceType(ua)
		deviceGroup := device.GetDeviceGroup(preciseType)

		// 将设备信息注入 Context
		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyDeviceID, deviceID)
		ctx = context.WithValue(ctx, KeyDeviceType, preciseType)
		ctx = context.WithValue(ctx, KeyDeviceGroup, deviceGroup)

		next(w, r.WithContext(ctx))
	}
}
