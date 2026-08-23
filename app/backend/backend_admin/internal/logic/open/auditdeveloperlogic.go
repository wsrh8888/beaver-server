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

package open

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AuditDeveloperLogic) AuditDeveloper(req *types.AuditDeveloperReq) (resp *types.AuditDeveloperRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	id, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil {
		return nil, errors.New("无效的申请ID")
	}

	_, err = l.svcCtx.OpenRpc.AuditDeveloper(l.ctx, &open_rpc.AuditDeveloperReq{
		Id:          id,
		Status:      int32(req.Status),
		AuditBy:     req.UserID,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		l.Errorf("审核开发者失败: %v", err)
		return nil, err
	}

	return &types.AuditDeveloperRes{}, nil
}
