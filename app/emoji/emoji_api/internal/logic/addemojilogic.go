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
	mqwsconst "beaver/common/const/mqwsconst"
	"context"
	"errors"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type AddEmojiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAddEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiLogic {
	return &AddEmojiLogic{
		ctx:    ctx,
		logger: beaverlog.New("add_emoji", ctx),
		svcCtx: svcCtx,
	}
}

func (l *AddEmojiLogic) AddEmoji(req *types.AddEmojiReq) (resp *types.AddEmojiRes, err error) {
	// 先按 FileKey 查重，已有则复用，不重复落库
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("file_key = ?", req.FileKey).First(&emoji).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 仅当不存在时创建新的 emoji
	if errors.Is(err, gorm.ErrRecordNotFound) {
		emojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
		if emojiVersion == -1 {
			l.logger.Error(model.LogMsg{Text: "生成表情版本号失败", Data: map[string]interface{}{"fileKey": req.FileKey}})
			return nil, errors.New("生成版本号失败")
		}

		emoji = emoji_models.Emoji{
			EmojiID: uuid.New().String(),
			FileKey: req.FileKey,
			Title:   req.Title,
			Version: emojiVersion,
			EmojiInfo: emoji_models.EmojiInfo{
				Width:  req.EmojiInfo.Width,
				Height: req.EmojiInfo.Height,
			},
		}

		if err := l.svcCtx.DB.Create(&emoji).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "添加表情失败", Data: map[string]interface{}{"fileKey": req.FileKey, "err": err.Error()}})
			return nil, err
		}
	}

	// 生成收藏版本号（按用户ID分区）
	collectVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_collect", "user_id", req.UserID)
	if collectVersion == -1 {
		l.logger.Error(model.LogMsg{Text: "生成收藏版本号失败", Data: map[string]interface{}{"userId": req.UserID}})
		return nil, errors.New("生成版本号失败")
	}

	// 添加表情并收藏
	favoriteEmoji := emoji_models.EmojiCollectEmoji{
		EmojiCollectID: uuid.New().String(),
		UserID:         req.UserID,
		EmojiID:        emoji.EmojiID,
		Version:        collectVersion,
		PackageID:      req.PackageID,
	}

	// 去重：同一用户对同一 emoji 已收藏则跳过创建
	var existFavorite emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ? AND emoji_id = ?", req.UserID, emoji.EmojiID).First(&existFavorite).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := l.svcCtx.DB.Create(&favoriteEmoji).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "收藏表情失败", Data: map[string]interface{}{"userId": req.UserID, "emojiId": emoji.EmojiID, "err": err.Error()}})
			return nil, err
		}

		// 异步通过 WS 通知用户其他客户端
		go func(etcdAddr string, userId string, emojiCollectId string, collectVersion int64, emojiVersion int64) {
			// 构建表更新数据
			var tableUpdates []map[string]interface{}

			// 通知表情表更新（如果创建了新表情）
			if emojiVersion > 0 {
				emojiUpdates := map[string]interface{}{
					"table": "emoji",
					"data": []map[string]interface{}{
						{
							"version": emojiVersion,
							"emojiId": emoji.EmojiID,
						},
					},
				}
				tableUpdates = append(tableUpdates, emojiUpdates)
			}

			// 通知表情收藏表更新
			collectUpdates := map[string]interface{}{
				"table":  "emoji_collect",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"version":        collectVersion,
						"emojiCollectId": emojiCollectId,
					},
				},
			}
			tableUpdates = append(tableUpdates, collectUpdates)

			// 通知给自己（用户ID作为接收者，空字符串作为发送者表示系统操作）
			payload := map[string]interface{}{
				"command":        "EMOJI",
				"targetId":       userId,
				"type":           "emoji_receive",
				"body":           map[string]interface{}{"tableUpdates": tableUpdates},
				"conversationId": "",
			}
			l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
		}(l.svcCtx.Config.Etcd, req.UserID, favoriteEmoji.EmojiCollectID, collectVersion, emoji.Version)
	}

	l.logger.Info(model.LogMsg{
		Text: "表情收藏成功",
		Data: map[string]interface{}{
			"userId":  req.UserID,
			"emojiId": emoji.EmojiID,
		},
	})

	return &types.AddEmojiRes{}, nil
}
