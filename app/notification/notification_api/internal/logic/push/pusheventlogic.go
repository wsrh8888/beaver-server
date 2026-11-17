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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	notification_rpc "beaver/app/notification/notification_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
)

type PushEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：投递通知事件
func NewPushEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushEventLogic {
	return &PushEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("PushEvent", ctx),
	}
}

// PushEvent 只做协议转换，业务规则全部留在 notification_rpc。
func (l *PushEventLogic) PushEvent(req *types.PushEventReq) (resp *types.PushEventRes, err error) {
	rpcReq := &notification_rpc.PushEventReq{EventType: req.EventType, Category: req.Category, FromUserId: req.FromUserID, TargetId: req.TargetID, TargetType: req.TargetType, PayloadJson: req.PayloadJSON, ToUserIds: req.ToUserIDs, DedupHash: req.DedupHash}
	rpcRes, err := l.svcCtx.NotifyRpc.PushEvent(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.PushEventRes{EventID: rpcRes.EventId, Version: rpcRes.Version}, nil
}
