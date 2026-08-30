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
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/middleware/ua"
	"beaver/utils/device"
	"beaver/utils/jwts"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type QrcodeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQrcodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeStatusLogic {
	return &QrcodeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QrcodeStatusLogic) QrcodeStatus(req *types.QrcodeStatusReq) (*types.QrcodeStatusRes, error) {
	key := fmt.Sprintf(QrcodeKeyFmt, req.Token)
	sessionStr, err := l.svcCtx.Redis.Get(key).Result()
	if err == redis.Nil {
		return &types.QrcodeStatusRes{Status: QrcodeStatusExpired}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}

	var session QrcodeSession
	if err = json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}
	if session.Status == QrcodeStatusPending {
		return &types.QrcodeStatusRes{Status: QrcodeStatusPending}, nil
	}
	if session.Status != QrcodeStatusConfirmed {
		return &types.QrcodeStatusRes{Status: QrcodeStatusExpired}, nil
	}

	infoRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: session.ScannedUserID})
	if err != nil || infoRes.UserInfo == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	user := infoRes.UserInfo

	profile := ua.Profile(l.ctx)
	preciseType := ua.DeviceType(l.ctx)
	deviceGroup := ua.DeviceGroup(l.ctx)

	jwtExpireHours := QrcodeTokenExpireHours
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		UserID:   user.UserId,
		DeviceID: req.DeviceID,
	}, l.svcCtx.Config.Auth.AccessSecret, jwtExpireHours)
	if err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}

	loginKey := fmt.Sprintf("user_authentication_session:%s:%s", user.UserId, deviceGroup)
	loginInfo, _ := json.Marshal(map[string]any{
		"token": token, "device_id": req.DeviceID, "device_type": preciseType,
		"device_group": deviceGroup, "login_time": time.Now().Format("2006-01-02 15:04:05"),
		"source": session.Source,
	})
	if err = l.svcCtx.Redis.Set(loginKey, string(loginInfo), time.Duration(jwtExpireHours)*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("服务内部异常")
	}

	var credential auth_models.AuthCredentialModel
	if l.svcCtx.DB.Take(&credential, "user_id = ?", user.UserId).Error == nil {
		now := time.Now()
		credential.LastLoginAt = &now
		credential.LoginCount++
		_ = l.svcCtx.DB.Save(&credential).Error
	}

	_ = device.UpsertOnLogin(l.svcCtx.DB, user.UserId, req.DeviceID, profile, req.ClientIP)

	l.svcCtx.Redis.Del(key)
	return &types.QrcodeStatusRes{
		Status: QrcodeStatusConfirmed,
		Token:  token,
		UserID: user.UserId,
		Source: session.Source,
	}, nil
}
