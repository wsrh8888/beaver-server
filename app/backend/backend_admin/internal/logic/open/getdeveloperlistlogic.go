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
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperListLogic {
	return &GetDeveloperListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetDeveloperListLogic) GetDeveloperList(req *types.GetDeveloperListReq) (resp *types.GetDeveloperListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListDevelopers(l.ctx, &open_rpc.ListDevelopersReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Status:   int32(req.Status),
	})
	if err != nil {
		l.Errorf("获取开发者列表失败: %v", err)
		return nil, err
	}

	list := make([]types.DeveloperInfo, 0, len(rpcRes.List))
	for _, dev := range rpcRes.List {
		list = append(list, types.DeveloperInfo{
			ID:          fmt.Sprintf("%d", dev.Id),
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
			CreatedAt:   dev.CreatedAt,
		})
	}

	return &types.GetDeveloperListRes{Total: rpcRes.Total, List: list}, nil
}
