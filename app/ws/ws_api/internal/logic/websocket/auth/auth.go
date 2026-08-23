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

package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"beaver/utils/jwts"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

// VerifyWsToken 验证 WS 连接鉴权：JWT 签名 + Redis 登录态
// platform: desktop, mobile, web 等
func VerifyWsToken(token, accessSecret, userID, platform string, rdb *redis.Client) error {
	if token == "" {
		return errors.New("token不能为空")
	}

	claims, err := jwts.ParseToken(token, accessSecret)
	if err != nil {
		return errors.New("无效的token")
	}

	if claims.UserID != userID {
		return errors.New("token与用户不匹配")
	}

	// 1. 如果没有传 platform，则退化为遍历校验（为了兼容性）
	platforms := []string{platform}
	if platform == "" {
		platforms = []string{"desktop", "mobile"}
	}

	// 2. 检查 Redis 登录态：精准校验对应平台的登录信息
	for _, dt := range platforms {
		key := fmt.Sprintf("user_authentication_session:%s:%s", userID, dt)
		val, err := rdb.Get(key).Result()
		if err != nil || val == "" {
			continue
		}
		var info map[string]interface{}
		if json.Unmarshal([]byte(val), &info) != nil {
			continue
		}
		if storedToken, _ := info["token"].(string); storedToken == token {
			// 增加设备标识符 (GUID) 强一致性校验
			// 确保当前连接的 Token 确实是给当前这台物理设备签发的
			if claims.DeviceID != "" {
				if storedDeviceID, ok := info["device_id"].(string); ok && storedDeviceID != claims.DeviceID {
					logx.Errorf("设备指纹不匹配: 用户 %s, Token绑定的设备: %s, 当前登录态设备: %s",
						userID, claims.DeviceID, storedDeviceID)
					return errors.New("设备标识符不匹配，疑似冒用")
				}
			}
			return nil
		}
	}

	return errors.New("登录已过期或在其他设备登录")
}
