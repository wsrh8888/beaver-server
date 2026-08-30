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
	"errors"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"
)

const (
	packageContentActionAdd    int32 = 1 // 向包内添加表情（必要时新建表情记录）
	packageContentActionRemove int32 = 2 // 从包内移除表情关联
)

type UpdateEmojiPackageContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEmojiPackageContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmojiPackageContentLogic {
	return &UpdateEmojiPackageContentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateEmojiPackageContentLogic) UpdateEmojiPackageContent(in *emoji_rpc.UpdateEmojiPackageContentReq) (*emoji_rpc.UpdateEmojiPackageContentRes, error) {
	switch in.Action {
	case packageContentActionAdd:
		return l.addEmoji(in)
	case packageContentActionRemove:
		return l.removeEmoji(in)
	default:
		return nil, errors.New("不支持的操作类型")
	}
}

func (l *UpdateEmojiPackageContentLogic) addEmoji(in *emoji_rpc.UpdateEmojiPackageContentReq) (*emoji_rpc.UpdateEmojiPackageContentRes, error) {
	var pkg emoji_models.EmojiPackage
	if err := l.svcCtx.DB.Where("package_id = ?", in.PackageId).First(&pkg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情包不存在")
		}
		return nil, err
	}

	emojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
	if emojiVersion == -1 {
		return nil, errors.New("获取版本号失败")
	}

	info := emoji_models.EmojiInfo{}
	if in.EmojiInfo != nil {
		info.Width = int(in.EmojiInfo.Width)
		info.Height = int(in.EmojiInfo.Height)
	}
	emoji := emoji_models.Emoji{
		EmojiID:   uuid.New().String(),
		FileKey:   in.FileKey,
		Title:     in.Title,
		EmojiInfo: info,
		Status:    1,
		Version:   emojiVersion,
	}
	if err := l.svcCtx.DB.Create(&emoji).Error; err != nil {
		return nil, err
	}

	var relations []emoji_models.EmojiPackageEmoji
	if err := l.svcCtx.DB.Where("package_id = ?", pkg.PackageID).Find(&relations).Error; err != nil {
		return nil, err
	}
	maxSort := 0
	for _, r := range relations {
		if r.SortOrder > maxSort {
			maxSort = r.SortOrder
		}
	}

	contentVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_package_emoji", "package_id", pkg.PackageID)
	if contentVersion == -1 {
		return nil, errors.New("获取版本号失败")
	}

	relation := emoji_models.EmojiPackageEmoji{
		RelationID: uuid.New().String(),
		PackageID:  pkg.PackageID,
		EmojiID:    emoji.EmojiID,
		SortOrder:  maxSort + 1,
		Version:    contentVersion,
	}
	if err := l.svcCtx.DB.Create(&relation).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.UpdateEmojiPackageContentRes{
		RelationId: relation.RelationID,
		EmojiId:    emoji.EmojiID,
	}, nil
}

func (l *UpdateEmojiPackageContentLogic) removeEmoji(in *emoji_rpc.UpdateEmojiPackageContentReq) (*emoji_rpc.UpdateEmojiPackageContentRes, error) {
	var relation emoji_models.EmojiPackageEmoji
	if err := l.svcCtx.DB.Where("package_id = ? AND emoji_id = ?", in.PackageId, in.EmojiId).
		First(&relation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情不在该表情包中")
		}
		return nil, err
	}

	if err := l.svcCtx.DB.Delete(&relation).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.UpdateEmojiPackageContentRes{EmojiId: in.EmojiId}, nil
}
