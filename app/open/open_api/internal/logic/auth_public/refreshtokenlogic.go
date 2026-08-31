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
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		logger: beaverlog.New("refresh_token", ctx),
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenRes, err error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refreshToken 不能为空")
	}

	var oldToken open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("refresh_token = ?", req.RefreshToken).First(&oldToken).Error; err != nil {
		l.logger.Warn(model.LogMsg{
			Text: "刷新令牌无效",
		})
		return nil, errors.New("刷新令牌无效")
	}
	if time.Now().Unix() > oldToken.RefreshTokenExpiresAt {
		l.logger.Warn(model.LogMsg{
			Text: "刷新令牌已过期",
			Data: map[string]any{"appId": oldToken.AppID},
		})
		return nil, errors.New("刷新令牌已过期，请重新授权")
	}

	rpcResp, err := l.svcCtx.OpenRpc.RefreshToken(l.ctx, &open_rpc.RefreshTokenReq{
		AppId:        oldToken.AppID,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "刷新令牌RPC失败",
			Data: map[string]any{"appId": oldToken.AppID, "err": err.Error()},
		})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "刷新令牌成功",
		Data: map[string]any{"appId": oldToken.AppID},
	})

	return &types.RefreshTokenRes{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
	}, nil
}
