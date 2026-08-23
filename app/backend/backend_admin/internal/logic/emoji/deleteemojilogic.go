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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiLogic {
	return &DeleteEmojiLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteEmoji 管理后台：删除表情。
// admin 职责：校验路径参数，将「删除」接口语义映射为 SaveEmoji.delete=true。
// RPC 职责：执行领域删除与版本同步，不与 HTTP 路由 1:1 命名。
func (l *DeleteEmojiLogic) DeleteEmoji(req *types.DeleteEmojiReq) (resp *types.DeleteEmojiRes, err error) {
	if req.EmojiId == "" {
		return nil, errors.New("表情ID不能为空")
	}

	del := true
	_, err = l.svcCtx.EmojiRpc.SaveEmoji(l.ctx, &emoji_rpc.SaveEmojiReq{
		EmojiId: req.EmojiId,
		Delete:  &del,
	})
	if err != nil {
		l.Errorf("删除表情失败: %v", err)
		return nil, err
	}
	return &types.DeleteEmojiRes{}, nil
}
