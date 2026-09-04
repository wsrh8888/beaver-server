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

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncMessageMediasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetSyncMessageMediasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMessageMediasLogic {
	return &GetSyncMessageMediasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_message_medias", ctx),
	}
}

func (l *GetSyncMessageMediasLogic) GetSyncMessageMedias(req *types.GetSyncMessageMediasReq) (*types.GetSyncMessageMediasRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.ChatRpc.GetSyncMessageMedias(l.ctx, &chat_rpc.GetSyncMessageMediasReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "同步消息媒体状态失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, errors.New("同步失败")
	}

	messageIDs := make([]string, 0)
	if rpcResp.MessageIds != nil {
		messageIDs = rpcResp.MessageIds
	}

	return &types.GetSyncMessageMediasRes{
		MessageIDs:      messageIDs,
		ServerTimestamp: rpcResp.ServerTimestamp,
	}, nil
}
