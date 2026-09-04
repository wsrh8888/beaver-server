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

type GetEmojiPackageCollectsByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 批量获取用户收藏的表情包记录详情（同步用）
func NewGetEmojiPackageCollectsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageCollectsByIdsLogic {
	return &GetEmojiPackageCollectsByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_emoji_package_collects_by_ids", ctx),
	}
}

func (l *GetEmojiPackageCollectsByIdsLogic) GetEmojiPackageCollectsByIds(req *types.GetEmojiPackageCollectsByIdsReq) (resp *types.GetEmojiPackageCollectsByIdsRes, err error) {
	if len(req.Ids) == 0 {
		return &types.GetEmojiPackageCollectsByIdsRes{
			Collects: make([]types.EmojiPackageCollectDetailItem, 0),
		}, nil
	}

	// 根据ID列表查询收藏记录详情
	var collects []emoji_models.EmojiPackageCollect
	err = l.svcCtx.DB.Where("package_collect_id IN ? AND user_id = ?", req.Ids, req.UserID).Find(&collects).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情包收藏记录详情失败", Data: map[string]any{"ids": req.Ids, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "批量查询表情包收藏记录详情", Data: map[string]interface{}{"requestCount": len(req.Ids), "resultCount": len(collects)}})

	// 转换为响应格式
	var collectItems []types.EmojiPackageCollectDetailItem
	for _, collect := range collects {
		collectItems = append(collectItems, types.EmojiPackageCollectDetailItem{
			PackageCollectID: collect.PackageCollectID,
			UserID:           collect.UserID,
			PackageID:        collect.PackageID,
			IsDeleted:        collect.IsDeleted,
			Version:          collect.Version,
			CreatedAt:        time.Time(collect.CreatedAt).UnixMilli(),
			UpdatedAt:        time.Time(collect.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetEmojiPackageCollectsByIdsRes{
		Collects: collectItems,
	}, nil
}
