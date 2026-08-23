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

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojisByUuidsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojisByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByUuidsLogic {
	return &GetEmojisByUuidsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojisByUuidsLogic) GetEmojisByUuids(req *types.GetEmojisByUuidsReq) (resp *types.GetEmojisByUuidsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojisByUuidsRes{Emojis: []types.EmojiSimpleItem{}}, nil
	}

	var emojis []emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id IN ? AND status = 1", req.Ids).Find(&emojis).Error; err != nil {
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
