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
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/auth/auth_api/internal/logic/auth_public"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/jwts"

	"github.com/go-redis/redis"
)

type QrcodeScanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewQrcodeScanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QrcodeScanLogic {
	return &QrcodeScanLogic{
		ctx:    ctx,
		logger: beaverlog.New("qrcode_scan", ctx),
		svcCtx: svcCtx,
	}
}

func (l *QrcodeScanLogic) QrcodeScan(req *types.QrcodeScanReq) (*types.QrcodeScanRes, error) {
	claims, err := jwts.ParseToken(req.AuthToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		l.logger.Warn(model.LogMsg{
			Text: "扫码鉴权失败",
			Data: map[string]any{"err": err.Error()},
		})
		return nil, fmt.Errorf("身份验证失败，请重新登录")
	}

	key := fmt.Sprintf(auth_public.QrcodeKeyFmt, req.Token)
	sessionStr, err := l.svcCtx.Redis.Get(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("二维码已过期，请刷新后重试")
	}
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "扫码读取会话失败",
			Data: map[string]any{"err": err.Error()},
		})
		return nil, fmt.Errorf("服务内部异常")
	}

	var session auth_public.QrcodeSession
	if err = json.Unmarshal([]byte(sessionStr), &session); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "扫码会话解析失败",
			Data: map[string]any{"err": err.Error()},
		})
		return nil, fmt.Errorf("服务内部异常")
	}
	if session.Status != auth_public.QrcodeStatusPending {
		return nil, fmt.Errorf("二维码已被使用或已过期")
	}

	ttl := l.svcCtx.Redis.TTL(key).Val()
	session.Status = auth_public.QrcodeStatusConfirmed
	session.ScannedUserID = claims.UserID
	updatedJSON, _ := json.Marshal(session)

	if err = l.svcCtx.Redis.Set(key, string(updatedJSON), ttl).Err(); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "扫码更新会话失败",
			Data: map[string]any{"err": err.Error()},
		})
		return nil, fmt.Errorf("服务内部异常")
	}

	l.logger.Info(model.LogMsg{
		Text: "扫码确认成功",
		Data: map[string]any{"userId": claims.UserID},
	})
	return &types.QrcodeScanRes{}, nil
}
