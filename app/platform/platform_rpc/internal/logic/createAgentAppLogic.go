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
	"crypto/rand"
	"encoding/hex"

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAgentAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentAppLogic {
	return &CreateAgentAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAgentAppLogic) CreateAgentApp(in *platform_rpc.CreateAgentAppReq) (*platform_rpc.CreateAgentAppRes, error) {
	appID := "app_" + randHex(16)
	row := agent_models.AgentApp{
		AppID:           appID,
		Name:            in.Name,
		Scene:           in.Scene,
		SystemPrompt:    in.SystemPrompt,
		LLMModelID:       in.LlmModelId,
		ModelParams:      in.ModelParams,
		Tools:           in.Tools,
		McpServerIDs:    in.McpServerIds,
		KnowledgeIDs:    in.KnowledgeIds,
		OpeningMessage:   in.OpeningMessage,
		OpeningQuestions: in.OpeningQuestions,
		Variables:        in.Variables,
		Status:          int(in.Status),
		Sort:            int(in.Sort),
		CreatedBy:       in.CreatedBy,
	}
	if row.Status == 0 {
		row.Status = 1
	}
	if row.CreatedBy == "" {
		row.CreatedBy = "system"
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		return &platform_rpc.CreateAgentAppRes{}, nil
	}
	return &platform_rpc.CreateAgentAppRes{AppId: appID}, nil
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
