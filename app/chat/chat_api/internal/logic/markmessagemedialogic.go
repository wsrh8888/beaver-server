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

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MarkMessageMediaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewMarkMessageMediaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkMessageMediaLogic {
	return &MarkMessageMediaLogic{
		ctx:    ctx,
		logger: logger.New("mark_message_media"),
		svcCtx: svcCtx,
	}
}

func (l *MarkMessageMediaLogic) MarkMessageMedia(req *types.MarkMessageMediaReq) (*types.MarkMessageMediaRes, error) {
	newMessageIDs := make([]string, 0, len(req.MessageIDs))

	for _, messageID := range req.MessageIDs {
		if messageID == "" {
			continue
		}

		var existing chat_models.ChatMessageMedia
		err := l.svcCtx.DB.Where("user_id = ? AND message_id = ?", req.UserID, messageID).First(&existing).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			logx.WithContext(l.ctx).Errorf("查询消息媒体状态失败: userId=%s, messageId=%s, error=%v", req.UserID, messageID, err)
			return nil, err
		}

		version := l.svcCtx.VersionGen.GetNextVersion("chat_message_medias", "user_id", req.UserID)
		record := chat_models.ChatMessageMedia{
			UserID:    req.UserID,
			MessageID: messageID,
			Version:   version,
		}
		if err := l.svcCtx.DB.Create(&record).Error; err != nil {
			logx.WithContext(l.ctx).Errorf("记录消息媒体状态失败: userId=%s, messageId=%s, error=%v", req.UserID, messageID, err)
			return nil, err
		}
		newMessageIDs = append(newMessageIDs, messageID)
	}

	if len(newMessageIDs) > 0 {
		go l.notifyMessageMedia(req.UserID, newMessageIDs)
	}

	l.logger.Info(model.LogMsg{
		Text: "标记消息媒体状态成功",
		Data: map[string]interface{}{
			"userId": req.UserID,
			"count":  len(newMessageIDs),
		},
	})

	return &types.MarkMessageMediaRes{}, nil
}

func (l *MarkMessageMediaLogic) notifyMessageMedia(userID string, messageIDs []string) {
	defer func() {
		if r := recover(); r != nil {
			logx.WithContext(l.ctx).Errorf("推送消息媒体状态失败 panic: %v", r)
		}
	}()

	update := map[string]interface{}{
		"table":  "message_medias",
		"userId": userID,
		"data": []map[string]interface{}{
			{"messageIds": messageIDs},
		},
	}

	payload := map[string]interface{}{
		"command":  wsCommandConst.CHAT_MESSAGE,
		"type":     wsTypeConst.ChatMessageMediaReceive,
		"senderId": userID,
		"targetId": userID,
		"body": map[string]interface{}{
			"tableUpdates": []map[string]interface{}{update},
		},
	}
	if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
		logx.WithContext(l.ctx).Errorf("推送消息媒体状态失败: userId=%s, error=%v", userID, err)
	}
}
