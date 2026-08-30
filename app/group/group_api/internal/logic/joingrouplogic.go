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
	"errors"
	"fmt"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		logger: logger.New("join_group"),
		svcCtx: svcCtx,
	}
}

func (l *JoinGroupLogic) JoinGroup(req *types.GroupJoinReq) (resp *types.GroupJoinRes, err error) {
	var memberVersion int64 = 0
	var invite *group_models.GroupInviteLinkModel

	groupID := req.GroupID
	if req.InviteCode != "" {
		var row group_models.GroupInviteLinkModel
		if e := l.svcCtx.DB.Where("token = ?", req.InviteCode).First(&row).Error; e != nil {
			return nil, fmt.Errorf("邀请无效")
		}
		if row.Status == 2 {
			return nil, fmt.Errorf("邀请已失效")
		}
		if row.Status == 3 || (row.MaxUses > 0 && row.UsedCount >= row.MaxUses) {
			return nil, fmt.Errorf("邀请已用尽")
		}
		if row.ExpireAt > 0 && time.Now().Unix() >= row.ExpireAt {
			return nil, fmt.Errorf("邀请已过期")
		}
		invite = &row
		if groupID != "" && groupID != invite.GroupID {
			return nil, fmt.Errorf("邀请与群不匹配")
		}
		groupID = invite.GroupID
		req.GroupID = groupID
	}
	if groupID == "" {
		return nil, fmt.Errorf("群ID不能为空")
	}

	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id = ? AND status = ?", groupID, 1).First(&group).Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("群组不存在或已解散，群组ID: %s", groupID)
		return nil, err
	}

	// 检查用户是否已经是群成员
	var existingMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ?", req.GroupID, req.UserID).First(&existingMember).Error
	if err == nil {
		// 用户已经是群成员
		if existingMember.Status == 1 {
			logx.WithContext(l.ctx).Errorf("用户已经是群成员，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
			return nil, err
		} else {
			// 用户之前被踢出，现在重新加入
			memberVersion = l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
			if memberVersion == -1 {
				logx.WithContext(l.ctx).Errorf("获取群成员版本号失败")
				return nil, errors.New("获取版本号失败")
			}
			err = l.svcCtx.DB.Model(&existingMember).Updates(map[string]interface{}{
				"status":    1,
				"join_time": time.Now(),
				"version":   memberVersion,
			}).Error
			if err != nil {
				logx.WithContext(l.ctx).Errorf("更新群成员状态失败: %v", err)
				return nil, err
			}

			// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel
		}
	} else {
		// 检查群组加入方式
		if group.JoinType == 1 {
			// 需要申请，创建申请记录
			// 获取该群入群申请的版本号（按群独立递增）
			requestVersion := l.svcCtx.VersionGen.GetNextVersion("group_join_requests", "group_id", req.GroupID)
			if requestVersion == -1 {
				logx.WithContext(l.ctx).Errorf("获取入群申请版本号失败")
				return nil, errors.New("获取版本号失败")
			}

			joinRequest := group_models.GroupJoinRequestModel{
				GroupID:         req.GroupID,
				ApplicantUserID: req.UserID,
				Message:         req.Message,
				Status:          0, // 待审核
				Version:         requestVersion,
			}
			err = l.svcCtx.DB.Create(&joinRequest).Error
			if err != nil {
				logx.WithContext(l.ctx).Errorf("创建入群申请失败: %v", err)
				return nil, err
			}

			resp = &types.GroupJoinRes{
				Version: requestVersion,
				Status:  0,
				GroupID: req.GroupID,
			}
			l.bumpInviteUse(invite)
			logx.WithContext(l.ctx).Infof("用户申请加入群组，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
			l.logger.Info(model.LogMsg{
				Text: "入群申请提交成功",
				Data: map[string]interface{}{
					"groupId": req.GroupID,
					"userId":  req.UserID,
				},
			})

			// 异步投递通知给群主/管理员
			go func() {
				ctx := context.Background()

				var admins []group_models.GroupMemberModel
				if err := l.svcCtx.DB.WithContext(ctx).
					Where("group_id = ? AND status = 1 AND role IN (?)", req.GroupID, []int{1, 2}).
					Find(&admins).Error; err != nil {
					logx.WithContext(l.ctx).Errorf("获取群管理员/群主失败(用于通知): %v", err)
					return
				}
				var toUsers []string
				for _, m := range admins {
					toUsers = append(toUsers, m.UserID)
				}
				if len(toUsers) == 0 {
					return
				}
				payload, _ := json.Marshal(map[string]interface{}{
					"requestId": requestVersion,
					"groupId":   req.GroupID,
					"userId":    req.UserID,
					"message":   req.Message,
				})
				_, err = l.svcCtx.NotifyRpc.PushEvent(ctx, &notification_rpc.PushEventReq{
					EventType:   notification_models.EventTypeGroupJoinRequest,
					Category:    notification_models.CategoryGroup,
					FromUserId:  req.UserID,
					TargetId:    req.GroupID,
					TargetType:  notification_models.TargetTypeGroup,
					PayloadJson: string(payload),
					ToUserIds:   toUsers,
					DedupHash:   fmt.Sprintf("%s_%d", req.GroupID, requestVersion),
				})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("投递入群申请通知失败: %v", err)
				}
			}()

			return resp, nil
		} else {
			// 获取该群成员的版本号（按群独立递增）
			memberVersion = l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
			if memberVersion == -1 {
				logx.WithContext(l.ctx).Errorf("获取群成员版本号失败")
				return nil, errors.New("获取版本号失败")
			}

			// 直接加入
			member := group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   req.UserID,
				Role:     3, // 普通成员
				Status:   1, // 正常状态
				JoinTime: time.Now(),
				Version:  memberVersion,
			}
			err = l.svcCtx.DB.Create(&member).Error
			if err != nil {
				logx.WithContext(l.ctx).Errorf("添加群成员失败: %v", err)
				return nil, err
			}

			// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel

			// 记录群成员变更日志
			changeLog := group_models.GroupMemberChangeLogModel{
				GroupID:    req.GroupID,
				UserID:     req.UserID,
				ChangeType: "join",
				OperatedBy: req.UserID,
				ChangeTime: time.Now(),
			}
			err = l.svcCtx.DB.Create(&changeLog).Error
			if err != nil {
				logx.WithContext(l.ctx).Errorf("记录群成员变更日志失败: %v", err)
				return nil, err
			}
		}
	}

	// 确保memberVersion有值（在直接加入的情况下）
	if memberVersion == 0 {
		memberVersion = l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
		if memberVersion == -1 {
			logx.WithContext(l.ctx).Errorf("获取群成员版本号失败")
			return nil, errors.New("获取版本号失败")
		}
	}

	// 更新新成员的会话记录（对标 groupmemberadd）
	_, err = l.svcCtx.ChatRpc.BatchUpdateConversation(l.ctx, &chat_rpc.BatchUpdateConversationReq{
		UserIds:        []string{req.UserID},
		ConversationId: "group_" + req.GroupID,
		LastMessage:    "",
	})
	if err != nil {
		logx.Errorf("Failed to update conversation: %v", err)
	}

	// 异步通知群成员
	go func() {
		ctx := context.Background()
		conversationID := "group_" + req.GroupID

		// 群聊内系统通知：xxx 加入了群聊（全员可见）
		if _, err := l.svcCtx.ChatRpc.SendNotificationMessage(ctx, &chat_rpc.SendNotificationMessageReq{
			ConversationId: conversationID,
			MessageType:    3,
			Content:        fmt.Sprintf("%s 加入了群聊", req.UserID),
			RelatedUserId:  req.UserID,
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("发送入群通知消息失败: %v", err)
		}

		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取群成员列表失败: %v", err)
			return
		}

		groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
		if groupVersion == -1 {
			logx.WithContext(l.ctx).Errorf("获取群组版本号失败")
		}

		joinMemberData := []map[string]interface{}{
			{
				"version": memberVersion,
				"groupId": req.GroupID,
				"userId":  req.UserID,
			},
		}

		// 1. 通知已在群的成员：group_members 变化
		for _, member := range response.Members {
			if member.UserID != req.UserID {
				payload := map[string]interface{}{
					"command":  wsCommandConst.GROUP_OPERATION,
					"type":     wsTypeConst.GroupMemberReceive,
					"senderId": req.UserID,
					"targetId": member.UserID,
					"body": map[string]interface{}{
						"tables": []map[string]interface{}{
							{
								"table": "group_members",
								"data":  joinMemberData,
							},
						},
					},
					"conversationId": "",
				}
				l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload)
			}
		}

		// 2. 通知新加入的成员：groups + group_members 变化
		joinerTables := []map[string]interface{}{
			{
				"table": "group_members",
				"data":  joinMemberData,
			},
		}
		if groupVersion != -1 {
			joinerTables = append([]map[string]interface{}{
				{
					"table": "groups",
					"data": []map[string]interface{}{
						{
							"version": groupVersion,
							"groupId": req.GroupID,
						},
					},
				},
			}, joinerTables...)
		}
		joinerPayload := map[string]interface{}{
			"command":  wsCommandConst.GROUP_OPERATION,
			"type":     wsTypeConst.GroupMemberReceive,
			"senderId": req.UserID,
			"targetId": req.UserID,
			"body": map[string]interface{}{
				"tables": joinerTables,
			},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, joinerPayload)

		// 3. 触发开放平台 Webhook 事件(群成员变更)
		l.triggerOpenPlatformWebhook(req.GroupID, req.UserID, []string{req.UserID}, "added")
	}()

	resp = &types.GroupJoinRes{
		Version: memberVersion,
		Status:  1,
		GroupID: req.GroupID,
	}

	l.bumpInviteUse(invite)

	logx.WithContext(l.ctx).Infof("用户加入群组成功，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "加入群组成功",
		Data: map[string]interface{}{
			"groupId": req.GroupID,
			"userId":  req.UserID,
		},
	})
	return resp, nil
}

func (l *JoinGroupLogic) bumpInviteUse(invite *group_models.GroupInviteLinkModel) {
	if invite == nil {
		return
	}
	updates := map[string]interface{}{
		"used_count": gorm.Expr("used_count + 1"),
	}
	if invite.MaxUses > 0 && invite.UsedCount+1 >= invite.MaxUses {
		updates["status"] = 3
	}
	_ = l.svcCtx.DB.Model(&group_models.GroupInviteLinkModel{}).Where("id = ?", invite.Id).Updates(updates).Error
}

func (l *JoinGroupLogic) triggerOpenPlatformWebhook(groupID string, operatorID string, memberIDs []string, action string) {
	defer func() {
		if r := recover(); r != nil {
			logx.WithContext(l.ctx).Errorf("触发开放平台 Webhook 时发生 panic: %v", r)
		}
	}()

	logx.WithContext(l.ctx).Infof("群成员变更事件: group_id=%s, action=%s, members=%v", groupID, action, memberIDs)
}
