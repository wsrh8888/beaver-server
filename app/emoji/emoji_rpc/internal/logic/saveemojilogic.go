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

type SaveEmojiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveEmojiLogic {
	return &SaveEmojiLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SaveEmojiLogic) SaveEmoji(in *emoji_rpc.SaveEmojiReq) (*emoji_rpc.SaveEmojiRes, error) {
	if in.Delete != nil && *in.Delete {
		return l.deleteEmoji(in.EmojiId)
	}
	if in.EmojiId == "" {
		return l.createEmoji(in)
	}
	return l.updateEmoji(in)
}

func (l *SaveEmojiLogic) createEmoji(in *emoji_rpc.SaveEmojiReq) (*emoji_rpc.SaveEmojiRes, error) {
	version := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
	if version == -1 {
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
		Version:   version,
	}
	if err := l.svcCtx.DB.Create(&emoji).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiRes{EmojiId: emoji.EmojiID}, nil
}

func (l *SaveEmojiLogic) updateEmoji(in *emoji_rpc.SaveEmojiReq) (*emoji_rpc.SaveEmojiRes, error) {
	var emoji emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id = ?", in.EmojiId).First(&emoji).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.PatchFileKey != nil {
		updates["file_key"] = *in.PatchFileKey
	}
	if in.PatchTitle != nil {
		if *in.PatchTitle != emoji.Title {
			var count int64
			if err := l.svcCtx.DB.Model(&emoji_models.Emoji{}).
				Where("title = ? AND emoji_id != ?", *in.PatchTitle, in.EmojiId).Count(&count).Error; err != nil {
				return nil, err
			}
			if count > 0 {
				return nil, errors.New("已存在同名表情")
			}
		}
		updates["title"] = *in.PatchTitle
	}
	if len(updates) == 0 {
		return &emoji_rpc.SaveEmojiRes{EmojiId: in.EmojiId}, nil
	}

	version := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
	if version == -1 {
		return nil, errors.New("获取版本号失败")
	}
	updates["version"] = version

	if err := l.svcCtx.DB.Model(&emoji).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiRes{EmojiId: in.EmojiId}, nil
}

func (l *SaveEmojiLogic) deleteEmoji(emojiID string) (*emoji_rpc.SaveEmojiRes, error) {
	var emoji emoji_models.Emoji
	if err := l.svcCtx.DB.Where("emoji_id = ?", emojiID).First(&emoji).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情不存在")
		}
		return nil, err
	}
	if err := l.svcCtx.DB.Delete(&emoji).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiRes{EmojiId: emojiID}, nil
}
