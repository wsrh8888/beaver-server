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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogsLogic {
	return &QueryLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogsLogic) QueryLogs(req *types.QueryLogsReq) (resp *types.QueryLogsRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.AdminQueryLogs(l.ctx, &platform_rpc.AdminQueryLogsReq{
		BucketId:   req.BucketID,
		Level:      req.Level,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Keyword:    req.Keyword,
		UserFilter: req.UserFilter,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("查询日志失败: %v", err)
		return nil, err
	}

	logs := make([]types.QueryLogsItem, 0, len(rpcRes.Logs))
	for _, item := range rpcRes.Logs {
		logs = append(logs, types.QueryLogsItem{
			Id:        uint(item.Id),
			Timestamp: item.Timestamp,
			Data:      item.Data,
		})
	}

	return &types.QueryLogsRes{
		Total: rpcRes.Total,
		Logs:  logs,
	}, nil
}
