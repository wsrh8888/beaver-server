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
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GroupInviteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 邀请用户加入群组
func NewGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInviteLogic {
	return &GroupInviteLogic{
		ctx:    ctx,
		logger: beaverlog.New("group_invite", ctx),
		svcCtx: svcCtx,
	}
}

func (l *GroupInviteLogic) GroupInvite(req *types.GroupInviteReq) (resp *types.GroupInviteRes, err error) {
	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id = ? AND status = ?", req.GroupID, 1).First(&group).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "群组不存在或已解散", Data: map[string]interface{}{"groupId": req.GroupID}})
		return nil, err
	}

	// 检查邀请者权限（群主或管理员）
	var inviterMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		req.GroupID, req.UserID, 1).First(&inviterMember).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "邀请者不是群成员", Data: map[string]interface{}{"groupId": req.GroupID, "userId": req.UserID}})
		return nil, err
	}

	// 检查邀请者角色（群主或管理员）
	if inviterMember.Role != 1 && inviterMember.Role != 2 {
		l.logger.Error(model.LogMsg{Text: "邀请者权限不足", Data: map[string]interface{}{"groupId": req.GroupID, "userId": req.UserID, "role": inviterMember.Role}})
		return nil, err
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// 处理每个被邀请的用户
	for _, userId := range req.UserIds {
		// 检查用户是否已经是群成员
		var existingMember group_models.GroupMemberModel
		err = tx.Where("group_id = ? AND user_id = ?", req.GroupID, userId).First(&existingMember).Error
		if err == nil {
			// 用户已经是群成员，更新状态为正常
			if existingMember.Status != 1 {
				err = tx.Model(&existingMember).Update("status", 1).Error
				if err != nil {
					tx.Rollback()
					l.logger.Error(model.LogMsg{Text: "更新群成员状态失败", Data: map[string]interface{}{"groupId": req.GroupID, "userId": userId, "err": err.Error()}})
					return nil, err
				}
			}
		} else {
			// 添加新群成员
			member := group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   userId,
				Role:     3, // 普通成员
				Status:   1, // 正常状态
				JoinTime: now,
				Version:  time.Now().Unix(),
			}
			err = tx.Create(&member).Error
			if err != nil {
				tx.Rollback()
				l.logger.Error(model.LogMsg{Text: "添加群成员失败", Data: map[string]interface{}{"groupId": req.GroupID, "userId": userId, "err": err.Error()}})
				return nil, err
			}
		}

		// 记录群成员变更日志
		changeLog := group_models.GroupMemberChangeLogModel{
			GroupID:    req.GroupID,
			UserID:     userId,
			ChangeType: "invite",
			OperatedBy: req.UserID,
			ChangeTime: now,
			Version:    time.Now().Unix(),
		}
		err = tx.Create(&changeLog).Error
		if err != nil {
			tx.Rollback()
			l.logger.Error(model.LogMsg{Text: "记录群成员变更日志失败", Data: map[string]interface{}{"groupId": req.GroupID, "userId": userId, "err": err.Error()}})
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "提交事务失败", Data: map[string]interface{}{"groupId": req.GroupID, "err": err.Error()}})
		return nil, err
	}

	// 获取该群成员的版本号（按群独立递增）
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
	if memberVersion == -1 {
		l.logger.Error(model.LogMsg{Text: "获取群成员版本号失败", Data: map[string]interface{}{"groupId": req.GroupID}})
		return nil, errors.New("获取版本号失败")
	}

	resp = &types.GroupInviteRes{
		Version: memberVersion,
	}

	// 异步通知相关成员
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表，用于推送通知
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "获取群成员列表失败", Data: map[string]interface{}{"groupId": req.GroupID, "err": err.Error()}})
			return
		}

		// 推送给已存在的群成员 - 群成员变动通知
		for _, member := range response.Members {
			if member.UserID != req.UserID { // 不通知操作者自己
				payload := map[string]interface{}{
					"command":  wsCommandConst.GROUP_OPERATION,
					"type":     wsTypeConst.GroupMemberReceive,
					"senderId": req.UserID,
					"targetId": member.UserID,
					"body": map[string]interface{}{
						"table": "group_members",
						"data": []map[string]interface{}{
							{
								"version": memberVersion,
								"groupId": req.GroupID,
							},
						},
					},
					"conversationId": "",
				}
				l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload)
			}
		}

		// 通知被邀请的成员
		for _, inviteeId := range req.UserIds {
			payload := map[string]interface{}{
				"command":  wsCommandConst.GROUP_OPERATION,
				"type":     wsTypeConst.GroupMemberReceive,
				"senderId": req.UserID,
				"targetId": inviteeId,
				"body": map[string]interface{}{
					"table": "group_members",
					"data": []map[string]interface{}{
						{
							"version": memberVersion,
							"groupId": req.GroupID,
						},
					},
				},
				"conversationId": "",
			}
			l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload)
		}
	}()

	l.logger.Info(model.LogMsg{
		Text: "群邀请成功",
		Data: map[string]interface{}{
			"groupId": req.GroupID,
			"userId":  req.UserID,
			"count":   len(req.UserIds),
		},
	})
	return resp, nil
}
