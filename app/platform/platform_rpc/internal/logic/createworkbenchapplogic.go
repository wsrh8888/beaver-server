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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkbenchAppLogic {
	return &CreateWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateWorkbenchAppLogic) CreateWorkbenchApp(in *platform_rpc.CreateWorkbenchAppReq) (*platform_rpc.CreateWorkbenchAppRes, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "应用名称不能为空")
	}

	appType := int8(in.AppType)
	if appType != 0 && appType != 1 {
		appType = 1
	}
	clientScope := int8(in.ClientScope)
	if clientScope != 0 && clientScope != 1 && clientScope != 2 {
		return nil, status.Error(codes.InvalidArgument, "可见端不合法")
	}

	entryConfig := fromProtoEntryConfig(in.EntryConfig)
	if msg := validateEntryConfig(appType, entryConfig); msg != "" {
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	statusVal := int8(in.Status)
	if statusVal != 0 && statusVal != 1 {
		return nil, status.Error(codes.InvalidArgument, "状态不合法")
	}
	openMode := int8(in.OpenMode)
	if openMode != 0 && openMode != 1 {
		return nil, status.Error(codes.InvalidArgument, "打开方式不合法")
	}

	app := platform_models.WorkbenchApp{
		WorkbenchAppID: strings.ReplaceAll(uuid.New().String(), "-", ""),
		Name:           name,
		Description:    strings.TrimSpace(in.Description),
		Icon:           strings.TrimSpace(in.Icon),
		AppType:        appType,
		ClientScope:    clientScope,
		EntryConfig:    entryConfig,
		OpenMode:       openMode,
		Category:       int8(in.Category),
		Sort:           int(in.Sort),
		Status:         statusVal,
		Remark:         strings.TrimSpace(in.Remark),
		CreatedBy:      in.OperatorId,
		LastModifiedBy: in.OperatorId,
	}
	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		l.Errorf("创建工作台应用失败: %v", err)
		return nil, status.Error(codes.Internal, "创建工作台应用失败")
	}

	return &platform_rpc.CreateWorkbenchAppRes{WorkbenchAppId: app.WorkbenchAppID}, nil
}
