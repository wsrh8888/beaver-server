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

package jwts

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 外部票受众。C 端用户、开放平台开发者、后台管理员三票同密钥同源，
// 只有 Audience 能区分平面；网关按路由组校验，杜绝拿普通用户票打后台接口。
const (
	AudienceUser      = "beaver:user"
	AudienceDeveloper = "beaver:developer"
	AudienceAdmin     = "beaver:admin"
)

type JwtPayLoad struct {
	NickName string `json:"nickName"`
	UserID   string `json:"userId"`
	DeviceID string `json:"deviceId"` // 绑定的物理设备唯一标识 (GUID)
}

type CustomClaims struct {
	JwtPayLoad
	jwt.RegisteredClaims
}

// GenToken 签发外部票。audience 取 AudienceUser / AudienceDeveloper / AudienceAdmin，
// 由签发票面所属平面决定，签错平面等于发错权限。
func GenToken(payload JwtPayLoad, accessSecret string, expires int, audience string) (string, error) {
	claims := CustomClaims{
		JwtPayLoad: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
			Audience:  jwt.ClaimStrings{audience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessSecret))
}

func ParseToken(tokenStr string, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if clains, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return clains, nil
	}
	return nil, err
}
