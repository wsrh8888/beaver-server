package device

import (
	"strings"
)

// 根据User-Agent识别设备类型
func GetDeviceType(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	// 识别桌面端
	if strings.Contains(userAgent, "beaver_desktop_windows") {
		return "desktop"
	} else if strings.Contains(userAgent, "beaver_desktop_macos") {
		return "desktop"
	} else if strings.Contains(userAgent, "beaver_desktop_linux") {
		return "desktop"
	}

	// 识别移动端
	if strings.Contains(userAgent, "beaver_mobile_android") {
		return "mobile"
	} else if strings.Contains(userAgent, "beaver_mobile_ios") {
		return "mobile"
	} else if strings.Contains(userAgent, "beaver_mobile_harmonyos") {
		return "mobile"
	}

	// 不是桌面端和移动端的，剩下的都是web端
	return "web"
}

// 验证设备ID格式
func IsValidDeviceID(deviceID string) bool {
	if deviceID == "" || len(deviceID) < 8 || len(deviceID) > 64 {
		return false
	}
	return true
}
