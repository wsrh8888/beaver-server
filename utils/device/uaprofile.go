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
	"net/url"
	"strings"
)

// UAProfile 从客户端 UA 解析的设备档案
type UAProfile struct {
	PreciseType string // windows / ios / android / macos / linux
	DeviceGroup string // desktop / mobile
	Model       string // iPhone17,3、SM-G991B、Windows-PC
	OsVersion   string // 18.2、10.0.19045
	DisplayName string // iPhone 17 Pro、DESKTOP-HOME
}

func ParseUAProfile(userAgent string) UAProfile {
	preciseType := GetDeviceType(userAgent)
	return UAProfile{
		PreciseType: preciseType,
		DeviceGroup: GetDeviceGroup(preciseType),
		Model:       parseUAToken(userAgent, "model"),
		OsVersion:   parseUAToken(userAgent, "os"),
		DisplayName: parseUAToken(userAgent, "name"),
	}
}

func parseUAToken(userAgent, key string) string {
	prefix := key + "/"
	idx := strings.Index(userAgent, prefix)
	if idx < 0 {
		return ""
	}
	rest := userAgent[idx+len(prefix):]
	end := strings.IndexAny(rest, " )")
	raw := rest
	if end >= 0 {
		raw = rest[:end]
	}
	decoded, _ := url.QueryUnescape(strings.TrimSpace(raw))
	return decoded
}
