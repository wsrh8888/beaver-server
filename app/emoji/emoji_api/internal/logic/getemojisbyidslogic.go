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

type GetEmojisByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情详情（用于数据同步）
func NewGetEmojisByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojisByIdsLogic {
	return &GetEmojisByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmojisByIdsLogic) GetEmojisByIds(req *types.GetEmojisByIdsReq) (resp *types.GetEmojisByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojisByIdsRes{
			Emojis: make([]types.EmojiDetailItem, 0),
		}, nil
	}

	// 根据ID列表查询表情详情
	var emojis []emoji_models.Emoji
	err = l.svcCtx.DB.Where("emoji_id IN ? AND status = ?", req.Ids, 1).Find(&emojis).Error
	if err != nil {
		l.Errorf("查询表情详情失败: ids=%v, error=%v", req.Ids, err)
		return nil, err
	}

	l.Infof("批量查询表情详情: 请求%d个, 返回%d个", len(req.Ids), len(emojis))

	// 获取表情ID列表，用于查询关联的包信息
	emojiIDs := make([]string, len(emojis))
	for i, emoji := range emojis {
		emojiIDs[i] = emoji.EmojiID
	}

	// 查询表情包关联信息
	var packageEmojis []emoji_models.EmojiPackageEmoji
	if len(emojiIDs) > 0 {
		l.svcCtx.DB.Where("emoji_id IN ?", emojiIDs).Find(&packageEmojis)
	}

	// 建立表情ID到包ID的映射
	emojiToPackage := make(map[string]*string)
	for _, pe := range packageEmojis {
		if pe.PackageID != "" {
			emojiToPackage[pe.EmojiID] = &pe.PackageID
		}
	}

	// 转换为响应格式
	var emojiItems []types.EmojiDetailItem
	for _, emoji := range emojis {
		emojiItems = append(emojiItems, types.EmojiDetailItem{
			EmojiID: emoji.EmojiID,
			FileKey: emoji.FileKey,
			Title:   emoji.Title,
			EmojiInfo: types.GetEmojiByIdsInfo{
				Width:  emoji.EmojiInfo.Width,
				Height: emoji.EmojiInfo.Height,
			},
			PackageID: emojiToPackage[emoji.EmojiID],
			Status:    emoji.Status,
			Version:   emoji.Version,
			CreatedAt: time.Time(emoji.CreatedAt).UnixMilli(),
			UpdatedAt: time.Time(emoji.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojisByIdsRes{
		Emojis: emojiItems,
	}, nil
}
