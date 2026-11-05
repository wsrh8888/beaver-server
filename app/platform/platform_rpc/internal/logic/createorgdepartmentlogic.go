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

	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type CreateOrgDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateOrgDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrgDepartmentLogic {
	return &CreateOrgDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateOrgDepartmentLogic", ctx),
	}
}

func (l *CreateOrgDepartmentLogic) CreateOrgDepartment(in *platform_rpc.CreateOrgDepartmentReq) (*platform_rpc.CreateOrgDepartmentRes, error) {
	// todo: add your logic here and delete this line

	return &platform_rpc.CreateOrgDepartmentRes{}, nil
}
