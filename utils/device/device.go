package device

import (
	"strings"
)

const (
	DeviceIOS     = "ios"
	DeviceAndroid = "android"
	DeviceWindows = "windows"
	DeviceMacOS   = "macos"
	DeviceLinux   = "linux"
	DeviceUnknown = "illegal" // 不再使用 unknown，改用 illegal 标识非法接入
)

// GetDeviceType 根据 User-Agent 识别精准 OS
func GetDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// 桌面端识别
	if strings.Contains(ua, "beaverdesktop") {
		if strings.Contains(ua, "windows") {
			return DeviceWindows
		}
		if strings.Contains(ua, "mac") {
			return DeviceMacOS
		}
		if strings.Contains(ua, "linux") {
			return DeviceLinux
		}
	}

	// 移动端识别
	if strings.Contains(ua, "beavermobile") {
		if strings.Contains(ua, "ios") {
			return DeviceIOS
		}
		if strings.Contains(ua, "android") {
			return DeviceAndroid
		}
	}

	return DeviceUnknown
}

// GetDeviceGroup 获取设备族群，用于同族互踢（大厂通用逻辑：1个手机 + 1个电脑在线）
func GetDeviceGroup(deviceType string) string {
	switch deviceType {
	case DeviceWindows, DeviceMacOS, DeviceLinux:
		return "desktop"
	case DeviceIOS, DeviceAndroid:
		return "mobile"
	default:
		return "unknown"
	}
}

// 验证设备ID格式
func IsValidDeviceID(deviceID string) bool {
	if deviceID == "" || len(deviceID) < 8 || len(deviceID) > 64 {
		return false
	}
	return true
}
