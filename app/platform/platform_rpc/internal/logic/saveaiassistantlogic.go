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

type SaveAiAssistantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveAiAssistantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAiAssistantLogic {
	return &SaveAiAssistantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SaveAiAssistantLogic", ctx),
	}
}

func (l *SaveAiAssistantLogic) SaveAiAssistant(in *platform_rpc.SaveAiAssistantReq) (*platform_rpc.SaveAiAssistantRes, error) {
	if !ai_models.ValidScene(in.Scene) {
		return nil, errors.New("场景无效")
	}
	status := int(in.Status)
	if status == 0 {
		status = 1
	}
	tools := in.Tools
	if tools == "" {
		tools = "[]"
	}
	mcp := in.McpServerIds
	if mcp == "" {
		mcp = "[]"
	}
	params := in.ModelParams
	if params == "" {
		params = "{}"
	}

	var row ai_models.AiAssistant
	err := l.svcCtx.DB.Where("scene = ?", in.Scene).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		row = ai_models.AiAssistant{
			Scene:          in.Scene,
			Name:           in.Name,
			SystemPrompt:   in.SystemPrompt,
			LLMModelID:     in.LlmModelId,
			ModelParams:    params,
			Tools:          tools,
			McpServerIDs:   mcp,
			OpeningMessage: in.OpeningMessage,
			Status:         status,
		}
		if err := l.svcCtx.DB.Create(&row).Error; err != nil {
			return nil, errors.New("保存 AI 助手失败")
		}
		return &platform_rpc.SaveAiAssistantRes{}, nil
	}
	if err != nil {
		return nil, err
	}
	updates := map[string]any{
		"name":            in.Name,
		"system_prompt":   in.SystemPrompt,
		"llm_model_id":    in.LlmModelId,
		"model_params":    params,
		"tools":           tools,
		"mcp_server_ids":  mcp,
		"opening_message": in.OpeningMessage,
		"status":          status,
	}
	if err := l.svcCtx.DB.Model(&row).Updates(updates).Error; err != nil {
		return nil, errors.New("保存 AI 助手失败")
	}
	return &platform_rpc.SaveAiAssistantRes{}, nil
}
