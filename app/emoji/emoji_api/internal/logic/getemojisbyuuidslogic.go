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

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetEmojisByUuidsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetEmojisByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByUuidsLogic {
	return &GetEmojisByUuidsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_emojis_by_uuids", ctx),
	}
}

func (l *GetEmojisByUuidsLogic) GetEmojisByUuids(req *types.GetEmojisByUuidsReq) (resp *types.GetEmojisByUuidsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojisByUuidsRes{Emojis: []types.EmojiSimpleItem{}}, nil
	}

	var emojis []emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id IN ? AND status = 1", req.Ids).Find(&emojis).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "按UUID查询表情失败", Data: map[string]any{"ids": req.Ids, "err": err.Error()}})
		return nil, err
	}

	items := make([]types.EmojiSimpleItem, 0, len(emojis))
	for _, e := range emojis {
		items = append(items, types.EmojiSimpleItem{
			EmojiID: e.EmojiID,
			FileKey: e.FileKey,
			Title:   e.Title,
			Version: e.Version,
			Status:  e.Status,
			EmojiInfo: types.GetEmojiByUuidsInfo{
				Width:  e.EmojiInfo.Width,
				Height: e.EmojiInfo.Height,
			},
		})
	}

	resp = &types.GetEmojisByUuidsRes{
		Emojis: items,
	}
	return resp, nil
}
