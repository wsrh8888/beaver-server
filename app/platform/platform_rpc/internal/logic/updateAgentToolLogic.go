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

type UpdateAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentToolLogic {
	return &UpdateAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateAgentToolLogic) UpdateAgentTool(in *platform_rpc.UpdateAgentToolReq) (*platform_rpc.UpdateAgentToolRes, error) {
	updates := map[string]any{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.ToolType != "" {
		updates["tool_type"] = in.ToolType
	}
	if in.Config != "" {
		updates["config"] = in.Config
	}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	if in.Version != "" {
		updates["version"] = in.Version
	}
	if in.SkillMd != "" {
		updates["skill_md"] = in.SkillMd
	}
	if in.FileName != "" {
		updates["file_name"] = in.FileName
	}
	if in.Status != 0 {
		updates["status"] = int(in.Status)
	}
	if in.Sort != 0 {
		updates["sort"] = int(in.Sort)
	}
	if len(updates) > 0 {
		l.svcCtx.DB.Model(&agent_models.AgentTool{}).Where("tool_id = ?", in.ToolId).Updates(updates)
	}
	return &platform_rpc.UpdateAgentToolRes{}, nil
}
