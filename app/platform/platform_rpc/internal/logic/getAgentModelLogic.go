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

type GetAgentModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentModelLogic {
	return &GetAgentModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— 模型目录（agent/ 通用）——
// source=custom：按 (model_id, user_id) 取用户自定义模型。
// source=official（默认）：按 model_id 取官方模型；"auto" 取 sort 最小启用官方模型。
// 取到模型行后 JOIN 其绑定的接入通道，解密真实 Key（revealKey=true，运行时专供 agent 调厂商）。
func (l *GetAgentModelLogic) GetAgentModel(in *platform_rpc.GetAgentModelReq) (*platform_rpc.GetAgentModelRes, error) {
	source := in.ModelSource
	if source == "" {
		source = "official"
	}
	var row agent_models.AgentModel
	found := false
	if source == "custom" {
		if err := l.svcCtx.DB.Where("source = 'custom' AND model_id = ? AND user_id = ? AND status = 1",
			in.ModelId, in.UserId).First(&row).Error; err == nil {
			found = true
		}
	} else if in.ModelId == "" || in.ModelId == "auto" {
		if err := l.svcCtx.DB.Where("source = 'official' AND status = 1").
			Order("sort ASC, id ASC").First(&row).Error; err == nil {
			found = true
		}
	} else {
		if err := l.svcCtx.DB.Where("source = 'official' AND model_id = ? AND status = 1",
			in.ModelId).First(&row).Error; err == nil {
			found = true
		}
	}
	if !found {
		return &platform_rpc.GetAgentModelRes{}, nil
	}
	cred := fetchCredential(l.svcCtx.DB, row.CredentialID)
	return &platform_rpc.GetAgentModelRes{Model: toAgentModelItem(&row, cred, l.svcCtx.Secret, true)}, nil
}
