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
	"strings"
	"time"

	ai_models "beaver/app/platform/platform_models/ai"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertWorkspaceLogic {
	return &UpsertWorkspaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpsertWorkspaceLogic) UpsertWorkspace(in *platform_rpc.UpsertWorkspaceReq) (*platform_rpc.UpsertWorkspaceRes, error) {
	uid := strings.TrimSpace(in.UserId)
	did := strings.TrimSpace(in.DeviceId)
	spaceID := strings.TrimSpace(in.WorkspaceId)
	p := strings.TrimSpace(in.Path)
	nm := strings.TrimSpace(in.Name)
	if uid == "" || did == "" || spaceID == "" || p == "" {
		return &platform_rpc.UpsertWorkspaceRes{}, nil
	}
	folder := nm
	if folder == "" {
		folder = dirName(p)
	}
	var row ai_models.AgentWorkspace
	if err := l.svcCtx.DB.Where("id = ?", spaceID).First(&row).Error; err == nil {
		if row.UserID != uid {
			return &platform_rpc.UpsertWorkspaceRes{}, nil
		}
		l.svcCtx.DB.Model(&ai_models.AgentWorkspace{}).Where("id = ?", spaceID).Updates(map[string]any{
			"device_id": did,
			"name":      folder,
			"path":      p,
			"version":   time.Now().UnixMilli(),
		})
		return &platform_rpc.UpsertWorkspaceRes{}, nil
	}
	row = ai_models.AgentWorkspace{
		ID:       spaceID,
		UserID:   uid,
		DeviceID: did,
		Name:     folder,
		Path:     p,
		Version:  time.Now().UnixMilli(),
	}
	l.svcCtx.DB.Create(&row)
	return &platform_rpc.UpsertWorkspaceRes{}, nil
}
