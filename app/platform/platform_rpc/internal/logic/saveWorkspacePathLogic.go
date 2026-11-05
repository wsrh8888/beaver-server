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

type SaveWorkspacePathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveWorkspacePathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWorkspacePathLogic {
	return &SaveWorkspacePathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveWorkspacePath 记下客户端报上来的目录，返回这次要用的路径。
//   - 没带 path 沿用库里已有
//   - 目录和设备号都没变不写
//   - 不是这个用户的空间不改
func (l *SaveWorkspacePathLogic) SaveWorkspacePath(in *platform_rpc.SaveWorkspacePathReq) (*platform_rpc.SaveWorkspacePathRes, error) {
	reported := strings.TrimSpace(in.Path)
	var row ai_models.AgentWorkspace
	if err := l.svcCtx.DB.Where("id = ?", in.WorkspaceId).First(&row).Error; err == nil {
		// 不是这个用户的空间不改
		if row.UserID != in.UserId {
			return &platform_rpc.SaveWorkspacePathRes{Path: ""}, nil
		}
		current := strings.TrimSpace(row.Path)
		did := strings.TrimSpace(in.DeviceId)
		pathSame := reported == "" || reported == current
		deviceSame := did == "" || did == strings.TrimSpace(row.DeviceID)
		if pathSame && deviceSame {
			if current != "" {
				return &platform_rpc.SaveWorkspacePathRes{Path: current}, nil
			}
			return &platform_rpc.SaveWorkspacePathRes{Path: reported}, nil
		}
		updates := map[string]any{}
		if reported != "" {
			updates["path"] = reported
			if strings.TrimSpace(row.Name) == "" {
				updates["name"] = dirName(reported)
			}
		}
		if did != "" {
			updates["device_id"] = did
		}
		updates["version"] = time.Now().UnixMilli()
		l.svcCtx.DB.Model(&ai_models.AgentWorkspace{}).Where("id = ?", in.WorkspaceId).Updates(updates)
		result := strings.TrimSpace(row.Path)
		if reported != "" {
			result = reported
		}
		return &platform_rpc.SaveWorkspacePathRes{Path: result}, nil
	}
	// 新建
	if reported == "" || strings.TrimSpace(in.DeviceId) == "" {
		return &platform_rpc.SaveWorkspacePathRes{Path: reported}, nil
	}
	row = ai_models.AgentWorkspace{
		ID:       in.WorkspaceId,
		UserID:   in.UserId,
		DeviceID: in.DeviceId,
		Name:     dirName(reported),
		Path:     reported,
		Version:  time.Now().UnixMilli(),
	}
	l.svcCtx.DB.Create(&row)
	return &platform_rpc.SaveWorkspacePathRes{Path: reported}, nil
}

func dirName(path string) string {
	folder := strings.TrimRight(strings.ReplaceAll(path, "\\", "/"), "/")
	if idx := strings.LastIndex(folder, "/"); idx >= 0 {
		folder = folder[idx+1:]
	}
	if folder == "" {
		folder = "空间"
	}
	return folder[:min(128, len(folder))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
