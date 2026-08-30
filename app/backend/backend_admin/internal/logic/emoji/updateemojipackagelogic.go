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

type UpdateEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmojiPackageLogic {
	return &UpdateEmojiPackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// UpdateEmojiPackage 管理后台：更新表情包。
// admin 职责：校验 packageId，将部分更新字段转为 patch 语义（运营可能只改状态/封面）。
// RPC 职责：SaveEmojiPackage 更新分支。
func (l *UpdateEmojiPackageLogic) UpdateEmojiPackage(req *types.UpdateEmojiPackageReq) (resp *types.UpdateEmojiPackageRes, err error) {
	if req.PackageId == "" {
		return nil, errors.New("表情包ID不能为空")
	}

	rpcReq := &emoji_rpc.SaveEmojiPackageReq{PackageId: req.PackageId}
	if req.Title != nil {
		rpcReq.PatchTitle = req.Title
	}
	if req.CoverFile != nil {
		rpcReq.PatchCoverFile = req.CoverFile
	}
	if req.Description != nil {
		rpcReq.PatchDescription = req.Description
	}
	if req.Type != nil {
		rpcReq.PatchType = req.Type
	}
	if req.Status != nil {
		status := int32(*req.Status)
		rpcReq.PatchStatus = &status
	}

	_, err = l.svcCtx.EmojiRpc.SaveEmojiPackage(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("更新表情包失败: %v", err)
		return nil, err
	}
	return &types.UpdateEmojiPackageRes{}, nil
}
