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

package agent

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetAgentModelLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentModelLogic {
	return &GetAgentModelLogic{logger: beaverlog.New("get_agent_model", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentModelLogic) GetAgentModel(req *types.GetAgentModelReq) (resp *types.GetAgentModelRes, err error) {
	source := req.Source
	if source == "" {
		source = "official"
	}
	rpcRes, err := l.svcCtx.PlatformRpc.GetAgentModel(l.ctx, &platform_rpc.GetAgentModelReq{
		ModelId:     req.ModelID,
		ModelSource: source,
		UserId:      req.UserID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取模型失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	if rpcRes.Model == nil || rpcRes.Model.ModelId == "" {
		return nil, errors.New("模型不存在")
	}
	return &types.GetAgentModelRes{Model: toAgentModelInfo(rpcRes.Model)}, nil
}
