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
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateOpenAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOpenAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOpenAppsLogic {
	return &UpdateOpenAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateOpenAppsLogic) UpdateOpenApps(in *open_rpc.UpdateOpenAppsReq) (*open_rpc.UpdateOpenAppsRes, error) {
	if len(in.AppIds) == 0 {
		return &open_rpc.UpdateOpenAppsRes{}, nil
	}

	now := time.Now()
	updates := map[string]interface{}{"updated_at": now, "last_modified_by": in.OperatorId}

	switch in.Action {
	case 1: // 审核通过
		updates["audit_status"] = 1
		updates["status"] = 1
		updates["audited_by"] = in.OperatorId
		updates["audited_at"] = now
	case 2: // 审核拒绝
		updates["audit_status"] = 2
		updates["audited_by"] = in.OperatorId
		updates["audited_at"] = now
	case 3: // 禁用
		updates["status"] = 2
	case 4: // 启用（已发布）
		updates["status"] = 1
	default:
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}

	result := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("app_id IN ?", in.AppIds).Updates(updates)
	if result.Error != nil {
		l.Errorf("更新应用状态失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "操作失败")
	}

	return &open_rpc.UpdateOpenAppsRes{AffectedCount: result.RowsAffected}, nil
}
