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

type GetEmojiPackageContentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetEmojiPackageContentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsLogic {
	return &GetEmojiPackageContentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_emoji_package_contents", ctx),
	}
}

func (l *GetEmojiPackageContentsLogic) GetEmojiPackageContents(in *pb.GetEmojiPackageContentsReq) (*pb.GetEmojiPackageContentsRes, error) {
	var packageContents []emoji_models.EmojiPackageEmoji

	// 查询所有表情包内容（去掉时间戳过滤，确保同步所有数据）
	err := l.svcCtx.DB.Find(&packageContents).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询表情包内容失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "查询表情包内容版本变更", Data: map[string]interface{}{"count": len(packageContents)}})

	// 转换为版本摘要格式
	var contentVersions []*pb.EmojiPackageContentVersionItem
	for _, content := range packageContents {
		contentVersions = append(contentVersions, &pb.EmojiPackageContentVersionItem{
			PackageId: content.PackageID,
			Version:   content.Version,
		})
	}

	return &pb.GetEmojiPackageContentsRes{
		EmojiPackageContentVersions: contentVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
