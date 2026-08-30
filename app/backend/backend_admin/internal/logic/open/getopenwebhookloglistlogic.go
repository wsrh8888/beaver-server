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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOpenWebhookLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOpenWebhookLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOpenWebhookLogListLogic {
	return &GetOpenWebhookLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOpenWebhookLogListLogic) GetOpenWebhookLogList(req *types.GetOpenWebhookLogListReq) (resp *types.GetOpenWebhookLogListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListWebhookLogs(l.ctx, &open_rpc.ListWebhookLogsReq{
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
		AppId:     req.AppID,
		EventType: req.EventType,
		Status:    int32(req.Status),
	})
	if err != nil {
		l.Errorf("查询 Webhook 日志失败: %v", err)
		return nil, err
	}

	list := make([]types.OpenWebhookLogInfo, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		list = append(list, types.OpenWebhookLogInfo{
			ID:           item.Id,
			AppID:        item.AppId,
			EventID:      item.EventId,
			EventType:    item.EventType,
			TargetURL:    item.TargetUrl,
			HTTPStatus:   int(item.HttpStatus),
			LatencyMs:    item.LatencyMs,
			RetryCount:   int(item.RetryCount),
			Status:       int(item.Status),
			ErrorMessage: item.ErrorMessage,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &types.GetOpenWebhookLogListRes{Total: rpcRes.Total, List: list}, nil
}
