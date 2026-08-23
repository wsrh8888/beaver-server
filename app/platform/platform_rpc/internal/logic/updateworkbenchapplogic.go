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
	"errors"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkbenchAppLogic {
	return &UpdateWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateWorkbenchAppLogic) UpdateWorkbenchApp(in *platform_rpc.UpdateWorkbenchAppReq) (*platform_rpc.UpdateWorkbenchAppRes, error) {
	if strings.TrimSpace(in.WorkbenchAppId) == "" {
		return nil, status.Error(codes.InvalidArgument, "应用 ID 不能为空")
	}

	var app platform_models.WorkbenchApp
	if err := l.svcCtx.DB.Where("workbench_app_id = ?", in.WorkbenchAppId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "工作台应用不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"last_modified_by": in.OperatorId,
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		updates["name"] = name
	}
	if in.Description != "" {
		updates["description"] = strings.TrimSpace(in.Description)
	}
	if in.Icon != "" {
		updates["icon"] = strings.TrimSpace(in.Icon)
	}
	if in.Remark != "" {
		updates["remark"] = strings.TrimSpace(in.Remark)
	}

	nextAppType := app.AppType
	if in.AppType != nil {
		appType := int8(*in.AppType)
		if appType != 0 && appType != 1 {
			return nil, status.Error(codes.InvalidArgument, "应用类型不合法")
		}
		updates["app_type"] = appType
		nextAppType = appType
	}
	if in.ClientScope != nil {
		clientScope := int8(*in.ClientScope)
		if clientScope != 0 && clientScope != 1 && clientScope != 2 {
			return nil, status.Error(codes.InvalidArgument, "可见端不合法")
		}
		updates["client_scope"] = clientScope
	}
	if in.EntryConfig != nil {
		entryConfig := fromProtoEntryConfig(in.EntryConfig)
		if msg := validateEntryConfig(nextAppType, entryConfig); msg != "" {
			return nil, status.Error(codes.InvalidArgument, msg)
		}
		updates["entry_config"] = entryConfig
	}
	if in.Category != nil {
		updates["category"] = int8(*in.Category)
	}
	if in.Sort != nil {
		updates["sort"] = int(*in.Sort)
	}
	if in.Status != nil {
		statusVal := int8(*in.Status)
		if statusVal != 0 && statusVal != 1 {
			return nil, status.Error(codes.InvalidArgument, "状态不合法")
		}
		updates["status"] = statusVal
	}
	if in.OpenMode != nil {
		openMode := int8(*in.OpenMode)
		if openMode != 0 && openMode != 1 {
			return nil, status.Error(codes.InvalidArgument, "打开方式不合法")
		}
		updates["open_mode"] = openMode
	}

	if err := l.svcCtx.DB.Model(&app).Updates(updates).Error; err != nil {
		l.Errorf("更新工作台应用失败: %v", err)
		return nil, status.Error(codes.Internal, "更新工作台应用失败")
	}

	return &platform_rpc.UpdateWorkbenchAppRes{}, nil
}
