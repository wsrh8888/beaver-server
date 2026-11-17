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

package admin

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	platformcli "beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
)

type AdminQueryLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：查询后台操作日志
func NewAdminQueryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminQueryLogsLogic {
	return &AdminQueryLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AdminQueryLogs", ctx),
	}
}

// AdminQueryLogs 只做协议转换，业务规则全部留在 platform_rpc。
func (l *AdminQueryLogsLogic) AdminQueryLogs(req *types.AdminQueryLogsReq) (resp *types.AdminQueryLogsRes, err error) {
	rpcReq := &platformcli.AdminQueryLogsReq{BucketId: req.BucketID, Level: req.Level, StartTime: req.StartTime, EndTime: req.EndTime, Keyword: req.Keyword, UserFilter: req.UserFilter, Page: req.Page, PageSize: req.PageSize}
	rpcRes, err := l.svcCtx.PlatformRpc.AdminQueryLogs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.AdminLogItem, 0, len(rpcRes.Logs))
	for _, it1 := range rpcRes.Logs {
		tmp2 = append(tmp2, types.AdminLogItem{ID: it1.Id, Level: it1.Level, Data: it1.Data, BucketID: it1.BucketId, BucketName: it1.BucketName, Timestamp: it1.Timestamp, CreatedAt: it1.CreatedAt})
	}
	return &types.AdminQueryLogsRes{Total: rpcRes.Total, Logs: tmp2}, nil
}
