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

package developer

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperDetailLogic {
	return &GetDeveloperDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDeveloperDetailLogic) GetDeveloperDetail(req *types.GetDeveloperDetailReq) (resp *types.GetDeveloperDetailRes, err error) {
	_, ok := l.ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("未登录")
	}
	if req.ID == 0 {
		return nil, errors.New("id 不能为空")
	}

	rpcRes, err := l.svcCtx.OpenRpc.GetDeveloper(l.ctx, &open_rpc.GetDeveloperReq{Id: uint64(req.ID)})
	if err != nil {
		return nil, err
	}

	dev := rpcRes.Developer
	return &types.GetDeveloperDetailRes{
		Developer: types.DeveloperInfo{
			ID:          uint(dev.Id),
			UserID:      dev.UserId,
			RealName:    dev.RealName,
			CompanyName: dev.CompanyName,
			Phone:       dev.Phone,
			Email:       dev.Email,
			Description: dev.Description,
			Status:      int(dev.Status),
			AuditBy:     dev.AuditBy,
			AuditTime:   dev.AuditTime,
			AuditRemark: dev.AuditRemark,
			CreatedAt:   dev.CreatedAt / 1000,
		},
	}, nil
}
