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

type DispatchPlatformEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDispatchPlatformEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchPlatformEventLogic {
	return &DispatchPlatformEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("dispatch_platform_event", ctx),
	}
}

// DispatchPlatformEvent 开放平台事件已改为通过长连接下发给机器人，
// 平台不再维护事件订阅表、也不做 Webhook 投递，因此这里不再有任何投递动作。
func (l *DispatchPlatformEventLogic) DispatchPlatformEvent(in *open_rpc.DispatchPlatformEventReq) (*open_rpc.DispatchPlatformEventRes, error) {
	return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
}
