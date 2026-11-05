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

	ai_models "beaver/app/platform/platform_models/ai"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLatestWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestWorkspaceLogic {
	return &GetLatestWorkspaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— 工作空间/设备（ai/ AI助手桌面端特有）——
// 取用户在某台电脑上最近同步过的工作空间。
func (l *GetLatestWorkspaceLogic) GetLatestWorkspace(in *platform_rpc.GetLatestWorkspaceReq) (*platform_rpc.GetLatestWorkspaceRes, error) {
	var row ai_models.AgentWorkspace
	if err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", in.UserId, in.DeviceId).
		Order("version DESC").First(&row).Error; err != nil {
		return &platform_rpc.GetLatestWorkspaceRes{}, nil
	}
	return &platform_rpc.GetLatestWorkspaceRes{
		Id:       row.ID,
		UserId:   row.UserID,
		DeviceId: row.DeviceID,
		Name:     row.Name,
		Path:     row.Path,
		Version:  row.Version,
	}, nil
}
