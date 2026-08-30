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
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/conversation"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type UpdateReadSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateReadSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReadSeqLogic {
	return &UpdateReadSeqLogic{
		ctx:    ctx,
		logger: beaverlog.New("update_read_seq", ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateReadSeqLogic) UpdateReadSeq(req *types.UpdateReadSeqReq) (*types.UpdateReadSeqRes, error) {
	if req.ConversationID == "" {
		l.logger.Warn(model.LogMsg{
			Text: "会话ID为空",
			Data: map[string]any{"userId": req.UserID},
		})
		return nil, errors.New("ConversationID不能为空")
	}
	if req.ReadSeq < 0 {
		l.logger.Warn(model.LogMsg{
			Text: "已读序列号非法",
			Data: map[string]any{
				"userId":         req.UserID,
				"conversationId": req.ConversationID,
				"readSeq":        req.ReadSeq,
			},
		})
		return nil, errors.New("ReadSeq值不对")
	}

	var userConvo chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, req.UserID).First(&userConvo).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)
			userConvo = chat_models.ChatUserConversation{
				UserID:         req.UserID,
				ConversationID: req.ConversationID,
				IsHidden:       false,
				IsPinned:       false,
				IsMuted:        false,
				UserReadSeq:    req.ReadSeq,
				Version:        version,
			}
			if err := l.svcCtx.DB.Create(&userConvo).Error; err != nil {
				l.logger.Error(model.LogMsg{
					Text: "创建用户会话关系失败",
					Data: map[string]any{
						"userId":         req.UserID,
						"conversationId": req.ConversationID,
						"err":            err.Error(),
					},
				})
				return nil, err
			}
			l.logger.Info(model.LogMsg{
				Text: "创建用户会话关系并设置已读序列号",
				Data: map[string]any{
					"userId":         req.UserID,
					"conversationId": req.ConversationID,
					"readSeq":        req.ReadSeq,
				},
			})
			go l.notifyReadSeqUpdate(req.UserID, req.ConversationID, version)
			go l.notifyPeerReadSeq(req.UserID, req.ConversationID, req.ReadSeq)
		} else {
			l.logger.Error(model.LogMsg{
				Text: "查询用户会话关系失败",
				Data: map[string]any{
					"userId":         req.UserID,
					"conversationId": req.ConversationID,
					"err":            err.Error(),
				},
			})
			return nil, err
		}
	} else {
		if req.ReadSeq > userConvo.UserReadSeq {
			version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)
			err = l.svcCtx.DB.Model(&userConvo).
				Updates(map[string]interface{}{
					"user_read_seq": req.ReadSeq,
					"updated_at":    time.Now(),
					"version":       version,
				}).Error
			if err != nil {
				l.logger.Error(model.LogMsg{
					Text: "更新已读序列号失败",
					Data: map[string]any{
						"userId":         req.UserID,
						"conversationId": req.ConversationID,
						"readSeq":        req.ReadSeq,
						"err":            err.Error(),
					},
				})
				return nil, err
			}
			l.logger.Info(model.LogMsg{
				Text: "更新已读序列号成功",
				Data: map[string]any{
					"userId":         req.UserID,
					"conversationId": req.ConversationID,
					"readSeq":        req.ReadSeq,
					"oldReadSeq":     userConvo.UserReadSeq,
				},
			})

			go l.notifyReadSeqUpdate(req.UserID, req.ConversationID, version)
			go l.notifyPeerReadSeq(req.UserID, req.ConversationID, req.ReadSeq)
		} else {
			l.logger.Info(model.LogMsg{
				Text: "已读序列号无需更新",
				Data: map[string]any{
					"userId":         req.UserID,
					"conversationId": req.ConversationID,
					"currentReadSeq": userConvo.UserReadSeq,
					"requestReadSeq": req.ReadSeq,
				},
			})
		}
	}

	return &types.UpdateReadSeqRes{
		Success: true,
	}, nil
}

func (l *UpdateReadSeqLogic) notifyReadSeqUpdate(userID, conversationID string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.logger.Error(model.LogMsg{
				Text: "推送已读同步异常",
				Data: map[string]any{
					"userId":         userID,
					"conversationId": conversationID,
					"panic":          fmt.Sprint(r),
				},
			})
		}
	}()

	userConversationsUpdate := map[string]interface{}{
		"table":          "user_conversations",
		"userId":         userID,
		"conversationId": conversationID,
		"data": []map[string]interface{}{
			{
				"version": version,
			},
		},
	}

	payload := map[string]interface{}{
		"command":  wsCommandConst.CHAT_MESSAGE,
		"type":     wsTypeConst.ChatUserConversationReceive,
		"senderId": userID,
		"targetId": userID,
		"body": map[string]interface{}{
			"tableUpdates": []map[string]interface{}{userConversationsUpdate},
		},
		"conversationId": conversationID,
	}
	if err := l.svcCtx.RocketMQ.SendMessage(l.ctx, mqwsconst.MqTopicWs, payload); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "推送已读同步失败",
			Data: map[string]any{
				"userId":         userID,
				"conversationId": conversationID,
				"err":            err.Error(),
			},
		})
		return
	}

	l.logger.Info(model.LogMsg{
		Text: "推送已读序列号更新通知完成",
		Data: map[string]any{
			"userId":         userID,
			"conversationId": conversationID,
			"version":        version,
		},
	})
}

func (l *UpdateReadSeqLogic) notifyPeerReadSeq(readerID, conversationID string, readSeq int64) {
	defer func() {
		if r := recover(); r != nil {
			l.logger.Error(model.LogMsg{
				Text: "推送对端已读异常",
				Data: map[string]any{
					"readerId":       readerID,
					"conversationId": conversationID,
					"panic":          fmt.Sprint(r),
				},
			})
		}
	}()

	peerIDs := l.getConversationPeerIDs(readerID, conversationID)
	for _, peerID := range peerIDs {
		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     wsTypeConst.ChatPeerReadReceive,
			"senderId": readerID,
			"targetId": peerID,
			"body": map[string]interface{}{
				"readerId":       readerID,
				"conversationId": conversationID,
				"readSeq":        readSeq,
			},
			"conversationId": conversationID,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
			l.logger.Error(model.LogMsg{
				Text: "推送对端已读通知失败",
				Data: map[string]any{
					"readerId":       readerID,
					"targetId":       peerID,
					"conversationId": conversationID,
					"err":            err.Error(),
				},
			})
		}
	}
}

func (l *UpdateReadSeqLogic) getConversationPeerIDs(currentUserID, conversationID string) []string {
	convType, userIDs := conversation.ParseConversationWithType(conversationID)
	if convType == 1 {
		var peers []string
		for _, uid := range userIDs {
			if uid != currentUserID {
				peers = append(peers, uid)
			}
		}
		return peers
	}

	var userConversations []chat_models.ChatUserConversation
	if err := l.svcCtx.DB.Where("conversation_id = ? AND user_id <> ?", conversationID, currentUserID).
		Find(&userConversations).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询群成员失败",
			Data: map[string]any{
				"conversationId": conversationID,
				"err":            err.Error(),
			},
		})
		return nil
	}

	peers := make([]string, 0, len(userConversations))
	for _, uc := range userConversations {
		peers = append(peers, uc.UserID)
	}
	return peers
}
