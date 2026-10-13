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

	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApplyDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("apply_developer", ctx)}
}

// ApplyDeveloper 已取消开发者表：登录用户即可用开放平台，本接口空成功（兼容旧客户端）。
func (l *ApplyDeveloperLogic) ApplyDeveloper(in *open_rpc.ApplyDeveloperReq) (*open_rpc.ApplyDeveloperRes, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "用户ID不能为空")
	}
	return &open_rpc.ApplyDeveloperRes{Id: 0}, nil
}
