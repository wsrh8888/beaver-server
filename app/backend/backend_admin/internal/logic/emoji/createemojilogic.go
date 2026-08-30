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

type CreateEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmojiLogic {
	return &CreateEmojiLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// CreateEmoji 管理后台：创建表情。
// admin 职责：校验运营录入的必填项（fileKey/title），将 HTTP 结构映射为领域写入请求。
// RPC 职责：SaveEmoji 创建分支（版本号、落库），可被其他服务复用。
func (l *CreateEmojiLogic) CreateEmoji(req *types.CreateEmojiReq) (resp *types.CreateEmojiRes, err error) {
	if req.FileUrl == "" {
		return nil, errors.New("文件地址不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("表情名称不能为空")
	}

	rpcRes, err := l.svcCtx.EmojiRpc.SaveEmoji(l.ctx, &emoji_rpc.SaveEmojiReq{
		FileKey: req.FileUrl,
		Title:   req.Title,
		EmojiInfo: &emoji_rpc.EmojiInfoMsg{
			Width:  int32(req.EmojiInfo.Width),
			Height: int32(req.EmojiInfo.Height),
		},
	})
	if err != nil {
		l.Errorf("创建表情失败: %v", err)
		return nil, err
	}
	return &types.CreateEmojiRes{EmojiId: rpcRes.EmojiId}, nil
}
