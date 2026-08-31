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
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取表情基础数据版本信息
func NewGetSyncEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojisLogic {
	return &GetSyncEmojisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_emojis", ctx),
	}
}

func (l *GetSyncEmojisLogic) GetSyncEmojis(req *types.GetSyncEmojisReq) (resp *types.GetSyncEmojisRes, err error) {
	// 调用Emoji RPC获取表情版本信息
	emojiResp, err := l.svcCtx.EmojiRpc.GetEmojis(l.ctx, &emoji_rpc.GetEmojisReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取表情版本信息失败", Data: map[string]any{"userId": req.UserID, "since": req.Since, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "查询表情版本信息完成",
		Data: map[string]interface{}{"userId": req.UserID, "count": len(emojiResp.EmojiVersions)},
	})

	// 转换为响应格式，确保返回空数组而不是null
	emojiVersions := make([]types.EmojiVersionItem, 0)
	if emojiResp.EmojiVersions != nil {
		for _, emoji := range emojiResp.EmojiVersions {
			emojiVersions = append(emojiVersions, types.EmojiVersionItem{
				EmojiId: emoji.EmojiId,
				Version: emoji.Version,
			})
		}
	}

	return &types.GetSyncEmojisRes{
		EmojiVersions:   emojiVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
