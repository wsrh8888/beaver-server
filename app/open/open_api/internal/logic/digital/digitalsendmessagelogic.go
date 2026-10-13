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

package digital

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	beaverlog "beaver/utils/beaverlog"
)

type DigitalSendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDigitalSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DigitalSendMessageLogic {
	return &DigitalSendMessageLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("digital_send_message", ctx)}
}

func (l *DigitalSendMessageLogic) DigitalSendMessage(req *types.DigitalSendMessageReq, authorization string) (resp *types.DigitalSendMessageRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppEnabled(app); err != nil {
		return nil, err
	}

	digital, err := utils.EnsureAppDigital(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	if req.ConversationID == "" || req.Content == "" {
		return nil, errors.New("conversationId 和 content 不能为空")
	}

	msgType := req.MsgType
	if msgType == "" {
		msgType = "text"
	}
	replyBody, err := buildDigitalChatMsg(msgType, req.Content)
	if err != nil {
		return nil, err
	}
	msg := replyBody
	if req.ReplyMsgID != "" {
		originSnap := &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: "[消息]"},
		}
		getRes, getErr := l.svcCtx.ChatRpc.GetChatMessage(l.ctx, &chat_rpc.GetChatMessageReq{
			ConversationId: req.ConversationID,
			MessageId:      req.ReplyMsgID,
		})
		if getErr == nil && getRes != nil && getRes.Found && getRes.Msg != nil {
			originSnap = getRes.Msg
		}
		msg = &chat_rpc.Msg{
			Type: 11,
			ReplyMsg: &chat_rpc.ReplyMsg{
				OriginMsgId: req.ReplyMsgID,
				OriginMsg:   originSnap,
				ReplyMsg:    replyBody,
			},
		}
	}

	messageID := fmt.Sprintf("digital_%s_%d", token.AppID, time.Now().UnixNano())
	chatRes, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:         digital.PartnerID,
		ConversationId: req.ConversationID,
		MessageId:      messageID,
		Msg:            msg,
		DeviceId:       "open_digital",
	})
	if err != nil {
		return nil, err
	}

	return &types.DigitalSendMessageRes{
		MessageID: chatRes.MessageId,
		SendTime:  time.Now().Unix(),
	}, nil
}

func buildDigitalChatMsg(msgType, content string) (*chat_rpc.Msg, error) {
	switch msgType {
	case "text", "":
		return &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: content},
		}, nil
	case "markdown":
		return &chat_rpc.Msg{
			Type: 13,
			MarkdownMsg: &chat_rpc.MarkdownMsg{
				Content: content,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported msgType: %s", msgType)
	}
}
