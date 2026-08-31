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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetEmojiPackageContentsByPackageIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 批量获取表情包内容详情（同步用）
func NewGetEmojiPackageContentsByPackageIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsByPackageIdsLogic {
	return &GetEmojiPackageContentsByPackageIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_emoji_package_contents_by_package_ids", ctx),
	}
}

func (l *GetEmojiPackageContentsByPackageIdsLogic) GetEmojiPackageContentsByPackageIds(req *types.GetEmojiPackageContentsByPackageIdsReq) (resp *types.GetEmojiPackageContentsByPackageIdsRes, err error) {
	if len(req.PackageIds) == 0 {
		return &types.GetEmojiPackageContentsByPackageIdsRes{
			Contents: make([]types.EmojiPackageContentDetailItem, 0),
		}, nil
	}

	// 根据表情包ID列表查询内容详情
	var contents []emoji_models.EmojiPackageEmoji
	err = l.svcCtx.DB.Where("package_id IN ?", req.PackageIds).Find(&contents).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情包内容详情失败", Data: map[string]any{"packageIds": req.PackageIds, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "批量查询表情包内容详情", Data: map[string]interface{}{"packageCount": len(req.PackageIds), "resultCount": len(contents)}})

	// 转换为响应格式
	var contentItems []types.EmojiPackageContentDetailItem
	for _, content := range contents {
		contentItems = append(contentItems, types.EmojiPackageContentDetailItem{
			RelationID: content.RelationID,
			PackageID:  content.PackageID,
			EmojiID:    content.EmojiID,
			SortOrder:  content.SortOrder,
			Version:    content.Version,
			CreatedAt:  time.Time(content.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(content.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiPackageContentsByPackageIdsRes{
		Contents: contentItems,
	}, nil
}
