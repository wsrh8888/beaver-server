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

const emojiPackageContentActionRemove int32 = 2

type RemoveEmojiFromPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveEmojiFromPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveEmojiFromPackageLogic {
	return &RemoveEmojiFromPackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// RemoveEmojiFromPackage 管理后台：从表情包移除表情。
// admin 职责：校验 packageId/emojiId，映射 action=移除。
// RPC 职责：UpdateEmojiPackageContent 维护包-表情关联，不与 HTTP 路由 1:1。
func (l *RemoveEmojiFromPackageLogic) RemoveEmojiFromPackage(req *types.RemoveEmojiFromPackageReq) (resp *types.RemoveEmojiFromPackageRes, err error) {
	if req.PackageId == "" {
		return nil, errors.New("表情包ID不能为空")
	}
	if req.EmojiId == "" {
		return nil, errors.New("表情ID不能为空")
	}

	_, err = l.svcCtx.EmojiRpc.UpdateEmojiPackageContent(l.ctx, &emoji_rpc.UpdateEmojiPackageContentReq{
		PackageId: req.PackageId,
		Action:    emojiPackageContentActionRemove,
		EmojiId:   req.EmojiId,
	})
	if err != nil {
		l.Errorf("从表情包移除表情失败: %v", err)
		return nil, err
	}
	return &types.RemoveEmojiFromPackageRes{}, nil
}
