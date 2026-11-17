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

	"beaver/app/asset/asset_api/internal/svc"
	"beaver/app/asset/asset_api/internal/types"
	assettypes "beaver/app/asset/asset_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
)

type EmojiCollectListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询用户收藏的表情
func NewEmojiCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmojiCollectListLogic {
	return &EmojiCollectListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("EmojiCollectList", ctx),
	}
}

// EmojiCollectList 只做协议转换，业务规则全部留在 asset_rpc。
func (l *EmojiCollectListLogic) EmojiCollectList(req *types.ListEmojiCollectsReq) (resp *types.ListEmojiCollectsRes, err error) {
	rpcReq := &assettypes.ListEmojiCollectsReq{Page: req.Page, PageSize: req.PageSize, UserId: req.UserID, EmojiId: req.EmojiID, StartTime: req.StartTime, EndTime: req.EndTime}
	rpcRes, err := l.svcCtx.EmojiRpc.ListEmojiCollects(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.EmojiCollectItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.EmojiCollectItem{CollectID: it1.CollectId, UserID: it1.UserId, EmojiID: it1.EmojiId, EmojiTitle: it1.EmojiTitle, EmojiFileKey: it1.EmojiFileKey, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListEmojiCollectsRes{Total: rpcRes.Total, List: tmp2}, nil
}
