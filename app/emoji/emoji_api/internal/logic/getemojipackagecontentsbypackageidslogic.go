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

type GetEmojiPackageContentsByPackageIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取表情包内容详情（同步用）
func NewGetEmojiPackageContentsByPackageIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsByPackageIdsLogic {
	return &GetEmojiPackageContentsByPackageIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
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
		l.Errorf("查询表情包内容详情失败: packageIds=%v, error=%v", req.PackageIds, err)
		return nil, err
	}

	l.Infof("批量查询表情包内容详情: 请求%d个包, 返回%d条内容", len(req.PackageIds), len(contents))

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
