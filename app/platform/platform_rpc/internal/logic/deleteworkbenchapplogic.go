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

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkbenchAppLogic {
	return &DeleteWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteWorkbenchAppLogic) DeleteWorkbenchApp(in *platform_rpc.DeleteWorkbenchAppReq) (*platform_rpc.DeleteWorkbenchAppRes, error) {
	if strings.TrimSpace(in.WorkbenchAppId) == "" {
		return nil, status.Error(codes.InvalidArgument, "应用 ID 不能为空")
	}

	result := l.svcCtx.DB.Where("workbench_app_id = ?", in.WorkbenchAppId).Delete(&platform_models.WorkbenchApp{})
	if result.Error != nil {
		l.Errorf("删除工作台应用失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "删除工作台应用失败")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "工作台应用不存在")
	}

	return &platform_rpc.DeleteWorkbenchAppRes{}, nil
}
