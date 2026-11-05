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

type GetAgentAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentAppLogic {
	return &GetAgentAppLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— AI 应用（agent/ 通用，渠道无关）——
// 按 app_id 取，或按 scene 取默认（sort 最小的启用项）。
func (l *GetAgentAppLogic) GetAgentApp(in *platform_rpc.GetAgentAppReq) (*platform_rpc.GetAgentAppRes, error) {
	var row agent_models.AgentApp
	if in.AppId != "" {
		if err := l.svcCtx.DB.Where("app_id = ? AND status = 1", in.AppId).First(&row).Error; err != nil {
			return &platform_rpc.GetAgentAppRes{}, nil
		}
	} else if in.Scene != "" {
		if err := l.svcCtx.DB.Where("scene = ? AND status = 1", in.Scene).
			Order("sort ASC, id ASC").First(&row).Error; err != nil {
			return &platform_rpc.GetAgentAppRes{}, nil
		}
	} else {
		return &platform_rpc.GetAgentAppRes{}, nil
	}
	return &platform_rpc.GetAgentAppRes{App: toAgentAppItem(&row)}, nil
}

func toAgentAppItem(row *agent_models.AgentApp) *platform_rpc.AgentAppItem {
	if row == nil {
		return nil
	}
	return &platform_rpc.AgentAppItem{
		Id:               uint64(row.Id),
		AppId:            row.AppID,
		Name:             row.Name,
		Scene:            row.Scene,
		SystemPrompt:     row.SystemPrompt,
		LlmModelId:       row.LLMModelID,
		ModelParams:      row.ModelParams,
		Tools:            row.Tools,
		McpServerIds:     row.McpServerIDs,
		KnowledgeIds:     row.KnowledgeIDs,
		OpeningMessage:   row.OpeningMessage,
		OpeningQuestions: row.OpeningQuestions,
		Variables:        row.Variables,
		Status:           int32(row.Status),
		Sort:             int32(row.Sort),
		CreatedBy:        row.CreatedBy,
	}
}
