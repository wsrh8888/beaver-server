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
	"strings"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ConversationInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConversationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationInfoLogic {
	return &ConversationInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConversationInfoLogic) ConversationInfo(req *types.ConversationInfoReq) (resp *types.ConversationInfoRes, err error) {
	var userConversation chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, req.UserID).First(&userConversation).Error

	resp = &types.ConversationInfoRes{
		ConversationID: req.ConversationID,
		MsgPreview:     "",
		UpdatedAt:      "",
		IsTop:          false,
	}

	var conversationMeta chat_models.ChatConversationMeta
	metaErr := l.svcCtx.DB.Where("conversation_id = ?", req.ConversationID).First(&conversationMeta).Error

	if err == nil {
		if metaErr == nil {
			resp.MsgPreview = conversationMeta.LastMessage
		}
		resp.UpdatedAt = userConversation.UpdatedAt.String()
		resp.IsTop = userConversation.IsPinned
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	convType := conversation.GetConversationType(req.ConversationID)
	switch convType {
	case 3: // 圈子（名称/头像由客户端补齐）
		resp.ChatType = 3
		resp.NickName = "圈子"
		resp.Avatar = ""
	case 1: // 私聊
		ids := strings.Split(req.ConversationID, "_")
		opponentID := ""
		if len(ids) >= 3 {
			if ids[1] == req.UserID {
				opponentID = ids[2]
			} else {
				opponentID = ids[1]
			}
		} else if len(ids) >= 2 {
			opponentID = ids[0]
			if ids[0] == req.UserID {
				opponentID = ids[1]
			}
		}
		userRes, userErr := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
			UserID: opponentID,
		})
		if userErr != nil {
			return nil, userErr
		}
		if userRes.UserInfo == nil {
			return nil, errors.New("user not found")
		}
		resp.Avatar = userRes.UserInfo.Avatar
		resp.NickName = userRes.UserInfo.NickName
		resp.ChatType = 1
	default: // 群聊
		groupID := req.ConversationID
		if strings.HasPrefix(groupID, "group_") {
			groupID = strings.TrimPrefix(groupID, "group_")
		}
		groupRes, groupErr := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
			GroupIDs: []string{groupID},
		})
		if groupErr != nil {
			return nil, groupErr
		}
		if len(groupRes.Groups) == 0 {
			return nil, errors.New("group not found")
		}
		group := groupRes.Groups[0]
		resp.Avatar = group.Avatar
		resp.NickName = group.Name
		resp.ChatType = 2
	}

	return resp, nil
}
