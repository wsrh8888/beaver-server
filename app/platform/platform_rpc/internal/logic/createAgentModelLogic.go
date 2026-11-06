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

type CreateAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentModelLogic {
	return &CreateAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentModelLogic) CreateAgentModel(in *platform_rpc.CreateAgentModelReq) (*platform_rpc.CreateAgentModelRes, error) {
	modelID := "amodel_" + randHex(16)
	source := in.Source
	if source == "" {
		source = "custom"
	}
	var userID *string
	if in.UserId != "" {
		uid := in.UserId
		userID = &uid
	}
	// 接入通道是厂商/地址/密钥的唯一来源：绑定不合法（不存在/停用/scope 不匹配）直接拒绝。
	if !validCredentialBinding(l.svcCtx.DB, in.CredentialId, source, in.UserId) {
		return &platform_rpc.CreateAgentModelRes{}, nil
	}
	var modelName *string
	if in.ModelName != "" {
		mn := in.ModelName
		modelName = &mn
	}
	status := int(in.Status)
	if status == 0 {
		status = 1
	}
	row := agent_models.AgentModel{
		ModelID:          modelID,
		Source:           source,
		UserID:           userID,
		Name:             in.Name,
		ModelName:        modelName,
		CredentialID:     in.CredentialId,
		ContextWindow:    in.ContextWindow,
		MaxOutputTokens:  int(in.MaxOutputTokens),
		InputModalities:  in.InputModalities,
		OutputModalities: in.OutputModalities,
		TokenTiers: &agent_models.TokenTiers{
			Flash:  int(in.MaxTokensFlash),
			Low:    int(in.MaxTokensLow),
			Medium: int(in.MaxTokensMedium),
			High:   int(in.MaxTokensHigh),
		},
		InstructionRole:  in.InstructionRole,
		SupportsThinking: in.SupportsThinking,
		ThinkingCanDisable: in.ThinkingCanDisable,
		ThinkingLevels:   in.ThinkingLevels,
		SupportsTools:    in.SupportsTools,
		SupportsVision:   in.SupportsVision,
		Sort:             int(in.Sort),
		Status:           status,
	}
	if row.InstructionRole == "" {
		row.InstructionRole = "system"
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentModelRes{}, nil
	}
	return &platform_rpc.CreateAgentModelRes{ModelId: modelID}, nil
}
