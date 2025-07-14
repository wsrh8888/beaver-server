package device

import (
	"strings"
)

// 根据User-Agent识别设备类型
func GetDeviceType(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	return "mobile"
	// 移动设备识别
	if strings.Contains(userAgent, "android") {
		return "mobile"
	} else if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") {
		return "mobile"
	} else if strings.Contains(userAgent, "mobile") {
		return "mobile"
	} else if strings.Contains(userAgent, "uniapp") {
		return "mobile"
	} else if strings.Contains(userAgent, "uni-app") {
		return "mobile"
	} else if strings.Contains(userAgent, "uni") {
		return "mobile"
	} else if strings.Contains(userAgent, "app") && (strings.Contains(userAgent, "android") || strings.Contains(userAgent, "ios")) {
		return "mobile"
	}

	// 桌面设备识别
	if strings.Contains(userAgent, "windows") {
		return "windows"
	} else if strings.Contains(userAgent, "macintosh") || strings.Contains(userAgent, "mac os") {
		return "mac"
	} else if strings.Contains(userAgent, "linux") {
		return "linux"
	} else {
		return "web"
	}
}

// 验证设备ID格式
func IsValidDeviceID(deviceID string) bool {
	// 设备ID不能为空且长度合理
	if deviceID == "" || len(deviceID) < 8 || len(deviceID) > 64 {
		return false
	}
	return true
}
