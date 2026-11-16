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

package emojilogic

import (
	"context"
	"time"

	emoji_models "beaver/app/asset/asset_models/emoji"
	"beaver/app/asset/asset_rpc/internal/svc"
	"beaver/app/asset/asset_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmojisLogic {
	return &ListEmojisLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_emojis", ctx)}
}

func (l *ListEmojisLogic) ListEmojis(in *pb.ListEmojisReq) (*pb.ListEmojisRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// package_id 非空：查包内表情，按 sort_order 排序
	if in.PackageId != "" {
		return l.listPackageEmojis(in.PackageId, page, pageSize)
	}

	db := l.svcCtx.DB.Model(&emoji_models.Emoji{})
	if in.EmojiId != "" {
		db = db.Where("emoji_id = ?", in.EmojiId)
	}
	if in.Title != "" {
		db = db.Where("title LIKE ?", "%"+in.Title+"%")
	}
	if in.StartTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.StartTime); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if in.EndTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", in.EndTime); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "统计表情失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}

	var list []emoji_models.Emoji
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情列表失败", Data: map[string]any{"page": page, "pageSize": pageSize, "err": err.Error()}})
		return nil, err
	}

	items := make([]*pb.EmojiItem, 0, len(list))
	for _, e := range list {
		items = append(items, toEmojiItem(e))
	}
	return &pb.ListEmojisRes{Total: total, List: items}, nil
}

func (l *ListEmojisLogic) listPackageEmojis(packageID string, page, pageSize int) (*pb.ListEmojisRes, error) {
	var relations []emoji_models.EmojiPackageEmoji
	if err := l.svcCtx.DB.Where("package_id = ?", packageID).
		Order("sort_order asc").Find(&relations).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情包内表情关联失败", Data: map[string]any{"packageId": packageID, "err": err.Error()}})
		return nil, err
	}
	if len(relations) == 0 {
		return &pb.ListEmojisRes{}, nil
	}

	emojiIDs := make([]string, len(relations))
	for i, r := range relations {
		emojiIDs[i] = r.EmojiID
	}

	var emojis []emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id IN ?", emojiIDs).Find(&emojis).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情包内表情详情失败", Data: map[string]any{"packageId": packageID, "err": err.Error()}})
		return nil, err
	}
	emojiMap := make(map[string]emoji_models.Emoji, len(emojis))
	for _, e := range emojis {
		emojiMap[e.EmojiID] = e
	}

	ordered := make([]emoji_models.Emoji, 0, len(relations))
	for _, r := range relations {
		if e, ok := emojiMap[r.EmojiID]; ok {
			ordered = append(ordered, e)
		}
	}

	total := int64(len(ordered))
	start := (page - 1) * pageSize
	if start >= len(ordered) {
		return &pb.ListEmojisRes{Total: total}, nil
	}
	end := start + pageSize
	if end > len(ordered) {
		end = len(ordered)
	}

	items := make([]*pb.EmojiItem, 0, end-start)
	for _, e := range ordered[start:end] {
		items = append(items, toEmojiItem(e))
	}
	return &pb.ListEmojisRes{Total: total, List: items}, nil
}

func toEmojiItem(e emoji_models.Emoji) *pb.EmojiItem {
	return &pb.EmojiItem{
		EmojiId:   e.EmojiID,
		FileKey:   e.FileKey,
		Title:     e.Title,
		EmojiInfo: &pb.EmojiInfoMsg{Width: int32(e.EmojiInfo.Width), Height: int32(e.EmojiInfo.Height)},
		Status:    int32(e.Status),
		CreatedAt: time.Time(e.CreatedAt).Format(time.RFC3339),
		UpdatedAt: time.Time(e.UpdatedAt).Format(time.RFC3339),
	}
}
