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

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetUserConversationSettingsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 批量获取用户会话设置数据
func NewGetUserConversationSettingsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationSettingsListByIdsLogic {
	return &GetUserConversationSettingsListByIdsLogic{
		ctx:    ctx,
		logger: beaverlog.New("get_user_conversation_settings", ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetUserConversationSettingsListByIdsLogic) GetUserConversationSettingsListByIds(req *types.GetUserConversationSettingsListByIdsReq) (resp *types.GetUserConversationSettingsListByIdsRes, err error) {
	var userConversations []chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("user_id = ? AND conversation_id IN (?)", req.UserID, req.ConversationIds).Find(&userConversations).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询用户会话设置失败",
			Data: map[string]any{
				"userId": req.UserID,
				"count":  len(req.ConversationIds),
				"err":    err.Error(),
			},
		})
		return nil, err
	}

	conversationSettings := make([]types.UserConversationSettingById, 0, len(userConversations))
	for _, uc := range userConversations {
		conversationSettings = append(conversationSettings, types.UserConversationSettingById{
			UserID:         uc.UserID,
			ConversationID: uc.ConversationID,
			IsHidden:       uc.IsHidden,
			IsPinned:       uc.IsPinned,
			IsMuted:        uc.IsMuted,
			UserReadSeq:    uc.UserReadSeq,
			Version:        uc.Version,
			CreatedAt:      time.Time(uc.CreatedAt).Unix(),
			UpdatedAt:      time.Time(uc.UpdatedAt).Unix(),
		})
	}

	l.logger.Info(model.LogMsg{
		Text: "批量查询用户会话设置成功",
		Data: map[string]any{
			"userId":       req.UserID,
			"requestCount": len(req.ConversationIds),
			"resultCount":  len(conversationSettings),
		},
	})
	return &types.GetUserConversationSettingsListByIdsRes{
		UserConversationSettings: conversationSettings,
	}, nil
}
