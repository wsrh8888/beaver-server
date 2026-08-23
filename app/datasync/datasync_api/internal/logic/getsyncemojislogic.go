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

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取表情基础数据版本信息
func NewGetSyncEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojisLogic {
	return &GetSyncEmojisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojisLogic) GetSyncEmojis(req *types.GetSyncEmojisReq) (resp *types.GetSyncEmojisRes, err error) {
	// 调用Emoji RPC获取表情版本信息
	emojiResp, err := l.svcCtx.EmojiRpc.GetEmojis(l.ctx, &emoji_rpc.GetEmojisReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情版本信息", len(emojiResp.EmojiVersions))

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
