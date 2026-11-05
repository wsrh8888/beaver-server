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

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentModelLogic {
	return &UpdateAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAgentModelLogic) UpdateAgentModel(in *platform_rpc.UpdateAgentModelReq) (*platform_rpc.UpdateAgentModelRes, error) {
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Provider != "" {
		updates["provider"] = in.Provider
	}
	if in.ModelName != "" {
		updates["model_name"] = in.ModelName
	}
	if in.Endpoint != "" {
		updates["endpoint"] = in.Endpoint
	}
	if in.ApiKey != "" {
		updates["api_key"] = in.ApiKey
	}
	updates["context_window"] = in.ContextWindow
	updates["max_output_tokens"] = int(in.MaxOutputTokens)
	if in.InputModalities != "" {
		updates["input_modalities"] = in.InputModalities
	}
	if in.OutputModalities != "" {
		updates["output_modalities"] = in.OutputModalities
	}
	if in.MaxTokensFlash != 0 || in.MaxTokensLow != 0 || in.MaxTokensMedium != 0 || in.MaxTokensHigh != 0 {
		updates["token_tiers"] = &agent_models.TokenTiers{
			Flash:  int(in.MaxTokensFlash),
			Low:    int(in.MaxTokensLow),
			Medium: int(in.MaxTokensMedium),
			High:   int(in.MaxTokensHigh),
		}
	}
	if in.InstructionRole != "" {
		updates["instruction_role"] = in.InstructionRole
	}
	if in.ThinkingLevels != "" {
		updates["thinking_levels"] = in.ThinkingLevels
	}
	updates["supports_tools"] = in.SupportsTools
	updates["supports_vision"] = in.SupportsVision
	updates["sort"] = int(in.Sort)
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if len(updates) > 0 {
		l.svcCtx.DB.Model(&agent_models.AgentModel{}).Where("model_id = ?", in.ModelId).Updates(updates)
	}
	return &platform_rpc.UpdateAgentModelRes{}, nil
}
