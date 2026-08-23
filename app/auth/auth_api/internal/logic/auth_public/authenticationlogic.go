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

package auth_public

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/middleware/ua"
	"beaver/utils/jwts"
	utils "beaver/utils/list"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationReq) (*types.AuthenticationRes, error) {
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		return &types.AuthenticationRes{}, nil
	}
	if req.Token == "" {
		return nil, errors.New("token不能为空")
	}

	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("认证失败")
	}

	deviceGroup, _ := l.ctx.Value(ua.KeyDeviceGroup).(string)
	groups := []string{deviceGroup}
	if deviceGroup == "" {
		groups = []string{"desktop", "mobile"}
	}

	var sessionKey string
	var loginInfo map[string]interface{}
	for _, g := range groups {
		if g == "" {
			continue
		}
		key := "user_authentication_session:" + claims.UserID + ":" + g
		val, err := l.svcCtx.Redis.Get(key).Result()
		if err != nil {
			continue
		}
		var info map[string]interface{}
		if json.Unmarshal([]byte(val), &info) != nil {
			continue
		}
		storedToken, _ := info["token"].(string)
		if storedToken != req.Token {
			continue
		}
		if claims.DeviceID != "" {
			if storedDeviceID, ok := info["device_id"].(string); ok && storedDeviceID != claims.DeviceID {
				return nil, errors.New("设备标识符不匹配")
			}
		}
		sessionKey = key
		loginInfo = info
		break
	}
	if loginInfo == nil {
		return nil, errors.New("token已失效")
	}

	loginInfo["last_active"] = time.Now().Format("2006-01-02 15:04:05")
	updated, _ := json.Marshal(loginInfo)
	l.svcCtx.Redis.Set(sessionKey, string(updated), time.Hour*time.Duration(l.svcCtx.Config.Auth.AccessExpire))

	return &types.AuthenticationRes{UserID: claims.UserID}, nil
}
