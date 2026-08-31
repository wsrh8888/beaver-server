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
	"encoding/json"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetForwardDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetForwardDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetForwardDetailsLogic {
	return &GetForwardDetailsLogic{
		ctx:    ctx,
		logger: beaverlog.New("get_forward_details", ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetForwardDetailsLogic) GetForwardDetails(req *types.GetForwardDetailsReq) (resp *types.GetForwardDetailsRes, err error) {
	var detail chat_models.ChatForward
	err = l.svcCtx.DB.Where("record_id = ?", req.RecordID).First(&detail).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "获取合并转发详情失败",
			Data: map[string]any{
				"recordId": req.RecordID,
				"err":      err.Error(),
			},
		})
		return nil, err
	}

	var userIds []string
	userIdSet := make(map[string]bool)
	for _, m := range detail.Content {
		if m.SendUserID != nil && *m.SendUserID != "" {
			if !userIdSet[*m.SendUserID] {
				userIds = append(userIds, *m.SendUserID)
				userIdSet[*m.SendUserID] = true
			}
		}
	}

	userInfoMap := make(map[string]types.Sender)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{
				Text: "批量获取用户信息失败",
				Data: map[string]any{
					"recordId": req.RecordID,
					"err":      err.Error(),
				},
			})
		} else {
			for userId, userInfo := range userListResp.UserInfo {
				userInfoMap[userId] = types.Sender{
					UserID:   userId,
					NickName: userInfo.NickName,
					Avatar:   userInfo.Avatar,
					UserType: int8(userInfo.UserType),
				}
			}
		}
	}

	var list []types.Message
	for _, m := range detail.Content {
		var tMsg types.Message
		tMsg.Id = m.Id
		tMsg.MessageID = m.MessageID
		tMsg.ConversationID = m.ConversationID
		tMsg.ConversationType = m.ConversationType
		tMsg.CreatedAt = m.CreatedAt.String()
		tMsg.Seq = m.Seq

		if m.Msg != nil {
			msgJSON, _ := json.Marshal(m.Msg)
			json.Unmarshal(msgJSON, &tMsg.Msg)
		}

		sendUserID := ""
		if m.SendUserID != nil {
			sendUserID = *m.SendUserID
		}

		if sendUserID != "" {
			if sender, exists := userInfoMap[sendUserID]; exists {
				tMsg.Sender = sender
			} else {
				tMsg.Sender = types.Sender{
					UserID:   sendUserID,
					NickName: "用户" + sendUserID[len(sendUserID)-4:],
					Avatar:   "",
				}
			}
		} else {
			tMsg.Sender = types.Sender{
				UserID:   "",
				NickName: "系统消息",
				Avatar:   "",
			}
		}

		list = append(list, tMsg)
	}

	l.logger.Info(model.LogMsg{
		Text: "获取合并转发详情成功",
		Data: map[string]any{
			"recordId": req.RecordID,
			"count":    len(list),
		},
	})
	return &types.GetForwardDetailsRes{
		List: list,
	}, nil
}
