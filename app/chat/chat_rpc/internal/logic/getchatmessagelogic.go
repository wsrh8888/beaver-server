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

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type GetChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageLogic {
	return &GetChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_chat_message", ctx),
	}
}

func (l *GetChatMessageLogic) GetChatMessage(in *chat_rpc.GetChatMessageReq) (*chat_rpc.GetChatMessageRes, error) {
	if in.ConversationId == "" || in.MessageId == "" {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	var row chat_models.ChatMessage
	err := l.svcCtx.DB.Where("conversation_id = ? AND message_id = ?", in.ConversationId, in.MessageId).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询聊天消息失败", Data: map[string]any{"conversationId": in.ConversationId, "messageId": in.MessageId, "err": err.Error()}})
		return nil, err
	}
	if row.Msg == nil {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	sl := NewSendMsgLogic(l.ctx, l.svcCtx)
	protoMsg, err := sl.convertCtypeMsgToGrpcMsg(*row.Msg)
	if err != nil || protoMsg == nil {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	return &chat_rpc.GetChatMessageRes{
		Found: true,
		Msg:   protoMsg,
	}, nil
}
