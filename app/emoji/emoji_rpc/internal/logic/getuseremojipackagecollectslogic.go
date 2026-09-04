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

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetUserEmojiPackageCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetUserEmojiPackageCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiPackageCollectsLogic {
	return &GetUserEmojiPackageCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_user_emoji_package_collects", ctx),
	}
}

func (l *GetUserEmojiPackageCollectsLogic) GetUserEmojiPackageCollects(in *emoji_rpc.GetUserEmojiPackageCollectsReq) (*emoji_rpc.GetUserEmojiPackageCollectsRes, error) {
	var packageCollects []emoji_models.EmojiPackageCollect
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&packageCollects).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询用户收藏表情包版本失败", Data: map[string]any{"userId": in.UserId, "since": in.Since, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "查询用户收藏表情包版本变更", Data: map[string]interface{}{"userId": in.UserId, "count": len(packageCollects)}})

	// 转换为版本摘要格式
	var packageCollectVersions []*emoji_rpc.EmojiPackageCollectVersionItem
	for _, pkgCollect := range packageCollects {
		packageCollectVersions = append(packageCollectVersions, &emoji_rpc.EmojiPackageCollectVersionItem{
			PackageCollectId: pkgCollect.PackageCollectID,
			Version:          pkgCollect.Version,
		})
	}

	return &emoji_rpc.GetUserEmojiPackageCollectsRes{
		EmojiPackageCollectVersions: packageCollectVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
