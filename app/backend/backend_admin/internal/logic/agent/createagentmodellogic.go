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

type CreateAgentModelLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentModelLogic {
	return &CreateAgentModelLogic{logger: beaverlog.New("create_agent_model", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateAgentModelLogic) CreateAgentModel(req *types.CreateAgentModelReq) (resp *types.CreateAgentModelRes, err error) {
	if req.Source != "official" && req.Source != "custom" {
		return nil, errors.New("来源只能是官方或个人")
	}
	if req.Source == "custom" && req.UserID == "" {
		return nil, errors.New("个人模型需要用户 ID")
	}
	rpcRes, err := l.svcCtx.PlatformRpc.CreateAgentModel(l.ctx, &platform_rpc.CreateAgentModelReq{
		Name:             req.Name,
		Source:           req.Source,
		UserId:           req.UserID,
		Provider:         req.Provider,
		ModelName:        req.ModelName,
		Endpoint:         req.Endpoint,
		ApiKey:           req.APIKey,
		ContextWindow:    req.ContextWindow,
		MaxOutputTokens:  int32(req.MaxOutputTokens),
		InputModalities:  "[]",
		OutputModalities: "[]",
		ThinkingLevels:   "[]",
		SupportsTools:    req.SupportsTools,
		SupportsVision:   req.SupportsVision,
		Sort:             int32(req.Sort),
		Status:           int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建模型失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	if rpcRes.ModelId == "" {
		return nil, errors.New("创建模型失败")
	}
	return &types.CreateAgentModelRes{ModelID: rpcRes.ModelId}, nil
}
