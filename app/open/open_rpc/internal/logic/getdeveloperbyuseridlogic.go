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
)

type GetDeveloperByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDeveloperByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperByUserIDLogic {
	return &GetDeveloperByUserIDLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("get_developer_by_user_id", ctx)}
}

// GetDeveloperByUserID 已取消开发者表：凡有 userId 即视为可用。
func (l *GetDeveloperByUserIDLogic) GetDeveloperByUserID(in *open_rpc.GetDeveloperByUserIDReq) (*open_rpc.GetDeveloperByUserIDRes, error) {
	if in.UserId == "" {
		return &open_rpc.GetDeveloperByUserIDRes{Found: false}, nil
	}
	return &open_rpc.GetDeveloperByUserIDRes{
		Found: true,
		Developer: &open_rpc.DeveloperItem{
			UserId: in.UserId,
			Status: 1,
		},
	}, nil
}
