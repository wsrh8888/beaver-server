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
	"fmt"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type InitializeConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

type initializedUserConversation struct {
	UserID  string
	Version int64
}

func NewInitializeConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitializeConversationLogic {
	return &InitializeConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("initialize_conversation", ctx),
	}
}

func (l *InitializeConversationLogic) InitializeConversation(in *chat_rpc.InitializeConversationReq) (*chat_rpc.InitializeConversationRes, error) {
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}
	if len(in.UserIds) == 0 {
		return nil, errors.New("用户列表不能为空")
	}
	if in.Type < 1 || in.Type > 4 {
		return nil, errors.New("无效的会话类型")
	}

	conversationVersion := int64(1)
	var existingConversation chat_models.ChatConversationMeta
	err := l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).First(&existingConversation).Error
	if err == nil {
		conversationVersion = existingConversation.Version
		if conversationVersion <= 0 {
			conversationVersion = 1
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		conversationMeta := chat_models.ChatConversationMeta{
			ConversationID: in.ConversationId,
			Type:           int(in.Type),
			MaxSeq:         0,
			Version:        1,
		}
		if err = l.svcCtx.DB.Create(&conversationMeta).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "创建会话元数据失败", Data: map[string]any{"conversationId": in.ConversationId, "err": err.Error()}})
			return nil, errors.New("创建会话元数据失败")
		}
	} else {
		l.logger.Error(model.LogMsg{Text: "查询会话失败", Data: map[string]any{"conversationId": in.ConversationId, "err": err.Error()}})
		return nil, errors.New("查询会话失败")
	}

	pushedUsers := make([]initializedUserConversation, 0, len(in.UserIds))

	for _, userId := range in.UserIds {
		var existingUC chat_models.ChatUserConversation
		ucErr := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", in.ConversationId, userId).
			First(&existingUC).Error
		if ucErr == nil {
			version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", userId)
			if err = l.svcCtx.DB.Model(&existingUC).Updates(map[string]interface{}{
				"is_hidden": false,
				"version":   version,
			}).Error; err != nil {
				l.logger.Error(model.LogMsg{Text: "更新用户会话关系失败", Data: map[string]any{"userId": userId, "conversationId": in.ConversationId, "err": err.Error()}})
				return nil, errors.New("更新用户会话关系失败")
			}
			pushedUsers = append(pushedUsers, initializedUserConversation{UserID: userId, Version: version})
			continue
		}
		if !errors.Is(ucErr, gorm.ErrRecordNotFound) {
			l.logger.Error(model.LogMsg{Text: "查询用户会话关系失败", Data: map[string]any{"userId": userId, "conversationId": in.ConversationId, "err": ucErr.Error()}})
			return nil, errors.New("查询用户会话关系失败")
		}

		userConversation := chat_models.ChatUserConversation{
			UserID:         userId,
			ConversationID: in.ConversationId,
			IsPinned:       false,
			IsMuted:        false,
			UserReadSeq:    0,
			Version:        1,
		}
		if err = l.svcCtx.DB.Create(&userConversation).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "创建用户会话关系失败", Data: map[string]any{"userId": userId, "conversationId": in.ConversationId, "err": err.Error()}})
			return nil, errors.New("创建用户会话关系失败")
		}
		pushedUsers = append(pushedUsers, initializedUserConversation{UserID: userId, Version: 1})
	}

	if len(pushedUsers) > 0 {
		go l.notifyConversationInitialized(in.ConversationId, conversationVersion, pushedUsers)
	}

	l.logger.Info(model.LogMsg{
		Text: "初始化会话成功",
		Data: map[string]interface{}{"conversationId": in.ConversationId, "type": in.Type, "userIds": in.UserIds},
	})

	return &chat_rpc.InitializeConversationRes{
		ConversationId: in.ConversationId,
	}, nil
}

func (l *InitializeConversationLogic) notifyConversationInitialized(
	conversationId string,
	conversationVersion int64,
	users []initializedUserConversation,
) {
	defer func() {
		if r := recover(); r != nil {
			l.logger.Error(model.LogMsg{Text: "推送会话初始化更新时发生panic", Data: map[string]any{"err": fmt.Sprintf("%v", r)}})
		}
	}()

	conversationsUpdate := map[string]interface{}{
		"table":          "conversations",
		"conversationId": conversationId,
		"data": []map[string]interface{}{
			{
				"version": int32(conversationVersion),
			},
		},
	}

	for _, user := range users {
		tableUpdates := []map[string]interface{}{
			conversationsUpdate,
			{
				"table":          "user_conversations",
				"userId":         user.UserID,
				"conversationId": conversationId,
				"data": []map[string]interface{}{
					{
						"version": int32(user.Version),
					},
				},
			},
		}

		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     wsTypeConst.ChatConversationMessageReceive,
			"senderId": "",
			"targetId": user.UserID,
			"body": map[string]interface{}{
				"tableUpdates": tableUpdates,
			},
			"conversationId": conversationId,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
			l.logger.Error(model.LogMsg{Text: "MQ推送会话初始化失败", Data: map[string]any{"recipient": user.UserID, "conversation": conversationId, "err": err.Error()}})
			continue
		}
		l.logger.Info(model.LogMsg{Text: "推送会话初始化成功", Data: map[string]interface{}{"recipient": user.UserID, "conversation": conversationId}})
	}
}
