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
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(req *types.ApplyDeveloperReq) (resp *types.ApplyDeveloperRes, err error) {
	userID := req.ApplicantUserID
	if userID == "" {
		userID = req.UserID
	}
	if userID == "" {
		return nil, errors.New("未指定申请用户")
	}

	rpcRes, err := l.svcCtx.OpenRpc.ApplyDeveloper(l.ctx, &open_rpc.ApplyDeveloperReq{
		UserId:      userID,
		RealName:    req.RealName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("开发者申请失败: %v", err)
		return nil, err
	}

	return &types.ApplyDeveloperRes{ID: fmt.Sprintf("%d", rpcRes.Id)}, nil
}
