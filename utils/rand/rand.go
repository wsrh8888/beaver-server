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

package utils

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

// 生成指定长度的随机字符串
func GenerateUUId() string {
	return uuid.New().String()
}

// 随机生成器中的字符集合
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// 生成指定长度的数字验证码
func GenerateNumericCode(length int) string {
	digits := "0123456789"
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}
	return sb.String()
}
