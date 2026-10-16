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

package logic

import (
	"context"
	"errors"

	ai_models "beaver/app/platform/platform_models/ai"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
	"gorm.io/gorm"
)

type GetAiAssistantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetAiAssistantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAiAssistantLogic {
	return &GetAiAssistantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetAiAssistantLogic", ctx),
	}
}

func (l *GetAiAssistantLogic) GetAiAssistant(in *platform_rpc.GetAiAssistantReq) (*platform_rpc.GetAiAssistantRes, error) {
	if !ai_models.ValidScene(in.Scene) {
		return nil, errors.New("场景无效")
	}
	var row ai_models.AiAssistant
	err := l.svcCtx.DB.Where("scene = ?", in.Scene).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &platform_rpc.GetAiAssistantRes{Assistant: &platform_rpc.AiAssistantItem{
			Scene:        in.Scene,
			Status:       1,
			Tools:        "[]",
			McpServerIds: "[]",
		}}, nil
	}
	if err != nil {
		return nil, err
	}
	return &platform_rpc.GetAiAssistantRes{Assistant: toAiAssistantItem(&row)}, nil
}

func toAiAssistantItem(row *ai_models.AiAssistant) *platform_rpc.AiAssistantItem {
	return &platform_rpc.AiAssistantItem{
		Scene:          row.Scene,
		Name:           row.Name,
		SystemPrompt:   row.SystemPrompt,
		LlmModelId:     row.LLMModelID,
		ModelParams:    row.ModelParams,
		Tools:          row.Tools,
		McpServerIds:   row.McpServerIDs,
		OpeningMessage: row.OpeningMessage,
		Status:         int32(row.Status),
	}
}
