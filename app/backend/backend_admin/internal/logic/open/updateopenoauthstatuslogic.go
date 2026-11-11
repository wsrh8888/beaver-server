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

package open

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateOpenOAuthStatusLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOpenOAuthStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOpenOAuthStatusLogic {
	return &UpdateOpenOAuthStatusLogic{logger: beaverlog.New("update_open_oauth_status", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateOpenOAuthStatusLogic) UpdateOpenOAuthStatus(req *types.UpdateOpenOAuthStatusReq) (resp *types.UpdateOpenOAuthStatusRes, err error) {
	if len(req.OAuthIDs) == 0 {
		return nil, errors.New("请选择应用")
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, errors.New("状态无效")
	}

	rpcRes, err := l.svcCtx.OpenRpc.UpdateOpenOAuth(l.ctx, &open_rpc.UpdateOpenOAuthReq{
		OauthIds:   req.OAuthIDs,
		Status:     int32(req.Status),
		OperatorId: req.UserID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "更新 OAuth 应用状态失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, err
	}

	return &types.UpdateOpenOAuthStatusRes{AffectedCount: rpcRes.AffectedCount}, nil
}
