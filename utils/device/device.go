/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

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
