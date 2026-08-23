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
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type DeleteRecentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewDeleteRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRecentLogic {
	return &DeleteRecentLogic{
		ctx:    ctx,
		logger: logger.New("delete_recent"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteRecentLogic) DeleteRecent(req *types.DeleteRecentReq) (resp *types.DeleteRecentRes, err error) {
	// 假删除操作
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{"is_delete": true, "is_pinned": false}).Error

	if err == nil {
		l.logger.Info(model.LogMsg{
			Text: "删除会话成功",
			Data: map[string]interface{}{
				"userId":         req.UserID,
				"conversationId": req.ConversationID,
			},
		})
	}

	return nil, err
}
