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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetAiAssistantLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAiAssistantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAiAssistantLogic {
	return &GetAiAssistantLogic{logger: beaverlog.New("get_ai_assistant", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAiAssistantLogic) GetAiAssistant(req *types.GetAiAssistantReq) (resp *types.GetAiAssistantRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetAiAssistant(l.ctx, &platform_rpc.GetAiAssistantReq{Scene: req.Scene})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取场景智能体失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	item := rpcRes.Assistant
	if item == nil {
		return &types.GetAiAssistantRes{Assistant: types.AiAssistantInfo{Scene: req.Scene, Status: 1, Tools: "[]", McpServerIDs: "[]"}}, nil
	}
	return &types.GetAiAssistantRes{Assistant: types.AiAssistantInfo{
		Scene:          item.Scene,
		Name:           item.Name,
		SystemPrompt:   item.SystemPrompt,
		LLMModelID:     item.LlmModelId,
		ModelParams:    item.ModelParams,
		Tools:          item.Tools,
		McpServerIDs:   item.McpServerIds,
		OpeningMessage: item.OpeningMessage,
		Status:         int(item.Status),
	}}, nil
}
