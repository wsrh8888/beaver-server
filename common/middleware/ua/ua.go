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
	KeyUAProfile   = "ua-profile"      // 完整设备档案（型号、系统版本、展示名）
)

func DeviceType(ctx context.Context) string {
	v, _ := ctx.Value(KeyDeviceType).(string)
	return v
}

func DeviceGroup(ctx context.Context) string {
	v, _ := ctx.Value(KeyDeviceGroup).(string)
	return v
}

func Profile(ctx context.Context) device.UAProfile {
	v, _ := ctx.Value(KeyUAProfile).(device.UAProfile)
	return v
}

// Middleware UA 识别中间件
// 从 User-Agent 中提取设备信息并注入 Context
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uaStr := r.Header.Get("User-Agent")
		deviceID := r.Header.Get("DeviceID")
		preciseType := device.GetDeviceType(uaStr)
		deviceGroup := device.GetDeviceGroup(preciseType)
		profile := device.ParseUAProfile(uaStr)

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyDeviceID, deviceID)
		ctx = context.WithValue(ctx, KeyDeviceType, preciseType)
		ctx = context.WithValue(ctx, KeyDeviceGroup, deviceGroup)
		ctx = context.WithValue(ctx, KeyUAProfile, profile)

		next(w, r.WithContext(ctx))
	}
}
