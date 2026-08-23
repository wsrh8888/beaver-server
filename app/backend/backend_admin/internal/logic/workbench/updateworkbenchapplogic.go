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

package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWorkbenchAppLogic {
	return &UpdateWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateWorkbenchAppLogic) UpdateWorkbenchApp(req *types.UpdateWorkbenchAppReq) (*types.UpdateWorkbenchAppRes, error) {
	in := &platform_rpc.UpdateWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
		Name:           req.Name,
		Description:    req.Description,
		Icon:           req.Icon,
		EntryConfig:    toProtoEntryConfig(req.EntryConfig),
		Remark:         req.Remark,
		OperatorId:     req.UserID,
	}
	if req.AppType != nil {
		v := int32(*req.AppType)
		in.AppType = &v
	}
	if req.ClientScope != nil {
		v := int32(*req.ClientScope)
		in.ClientScope = &v
	}
	if req.Category != nil {
		v := int32(*req.Category)
		in.Category = &v
	}
	if req.Sort != nil {
		sortVal := int32(*req.Sort)
		in.Sort = &sortVal
	}
	if req.Status != nil {
		statusVal := int32(*req.Status)
		in.Status = &statusVal
	}
	if req.OpenMode != nil {
		openMode := int32(*req.OpenMode)
		in.OpenMode = &openMode
	}

	_, err := l.svcCtx.PlatformRpc.UpdateWorkbenchApp(l.ctx, in)
	if err != nil {
		l.Errorf("更新工作台应用失败: %v", err)
		return nil, err
	}

	return &types.UpdateWorkbenchAppRes{}, nil
}
