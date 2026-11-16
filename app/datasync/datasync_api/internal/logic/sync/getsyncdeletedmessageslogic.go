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
	"errors"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	chat_rpc "beaver/app/message/message_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncDeletedMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取所有需要更新的已删除消息ID列表
func NewGetSyncDeletedMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncDeletedMessagesLogic {
	return &GetSyncDeletedMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_deleted_messages", ctx),
	}
}

func (l *GetSyncDeletedMessagesLogic) GetSyncDeletedMessages(req *types.GetSyncDeletedMessagesReq) (resp *types.GetSyncDeletedMessagesRes, err error) {
	// 删除记录属 message 库，统一走 message_rpc，不再直连其表
	rpcResp, err := l.svcCtx.MessageRpc.GetSyncDeletedMessages(l.ctx, &chat_rpc.GetSyncDeletedMessagesReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "同步已删除消息失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, errors.New("同步失败")
	}

	return &types.GetSyncDeletedMessagesRes{
		MessageIDs:      rpcResp.MessageIds,
		ServerTimestamp: rpcResp.ServerTimestamp,
	}, nil
}
