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

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectsByUuidsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按收藏ID批量获取表情收藏记录（同步补齐）
func NewGetEmojiCollectsByUuidsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectsByUuidsLogic {
	return &GetEmojiCollectsByUuidsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojiCollectsByUuidsLogic) GetEmojiCollectsByUuids(req *types.GetEmojiCollectsByUuidsReq) (resp *types.GetEmojiCollectsByUuidsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiCollectsByUuidsRes{
			Collects: make([]types.EmojiCollectDetailItem, 0),
		}, nil
	}

	var collects []emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("emoji_collect_id IN ?", req.Ids).Find(&collects).Error
	if err != nil {
		l.Errorf("按收藏ID批量查询表情收藏记录失败: ids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	collectItems := make([]types.EmojiCollectDetailItem, 0, len(collects))
	for _, collect := range collects {
		collectItems = append(collectItems, types.EmojiCollectDetailItem{
			EmojiCollectID: collect.EmojiCollectID,
			UserID:         collect.UserID,
			EmojiID:        collect.EmojiID,
			IsDeleted:      collect.IsDeleted,
			Version:        collect.Version,
			CreatedAt:      time.Time(collect.CreatedAt).UnixMilli(),
			UpdatedAt:      time.Time(collect.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiCollectsByUuidsRes{
		Collects: collectItems,
	}, nil
}
