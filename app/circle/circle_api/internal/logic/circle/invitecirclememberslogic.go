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

package circle

import (
	"context"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type InviteCircleMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewInviteCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteCircleMembersLogic {
	return &InviteCircleMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("invite_circle_members", ctx),
	}
}

func (l *InviteCircleMembersLogic) InviteCircleMembers(req *types.InviteCircleMembersReq) (resp *types.InviteCircleMembersRes, err error) {
	if len(req.UserIds) == 0 {
		return nil, fmt.Errorf("请选择要邀请的成员")
	}

	var circle circle_models.CircleModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可邀请成员")
	}

	addedUserIds := make([]string, 0, len(req.UserIds))
	for _, userID := range req.UserIds {
		if userID == "" || userID == req.UserID {
			continue
		}
		var existing circle_models.CircleMemberModel
		if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, userID).First(&existing).Error == nil {
			continue
		}
		member := circle_models.CircleMemberModel{
			CircleID: req.CircleID,
			UserID:   userID,
			Role:     3,
		}
		if err = l.svcCtx.DB.Create(&member).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "添加圈子成员失败", Data: map[string]any{"circleId": req.CircleID, "userId": userID, "err": err.Error()}})
			return nil, fmt.Errorf("添加成员失败: %v", err)
		}
		addedUserIds = append(addedUserIds, userID)
	}

	if len(addedUserIds) == 0 {
		return &types.InviteCircleMembersRes{}, nil
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Update("version", circleVersion)

	conversationID := fmt.Sprintf("circle_%s", req.CircleID)
	go func() {
		ctx := context.Background()
		l.svcCtx.ChatRpc.InitializeConversation(ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        addedUserIds,
		})

		for _, targetID := range addedUserIds {
			payload := map[string]interface{}{
				"command":  wsCommandConst.CIRCLE_OPERATION,
				"type":     wsTypeConst.CircleReceive,
				"senderId": req.UserID,
				"targetId": targetID,
				"body": map[string]interface{}{
					"tables": []map[string]interface{}{
						{
							"table": "circles",
							"data": []map[string]interface{}{
								{
									"version":  circleVersion,
									"circleId": req.CircleID,
								},
							},
						},
					},
				},
				"conversationId": conversationID,
			}
			if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
				l.logger.Error(model.LogMsg{Text: "推送圈子资料同步失败", Data: map[string]any{"circleId": req.CircleID, "targetId": targetID, "err": err.Error()}})
			}
		}
	}()

	l.logger.Info(model.LogMsg{Text: "邀请圈子成员成功", Data: map[string]interface{}{"circleId": req.CircleID, "addedUserIds": addedUserIds}})

	return &types.InviteCircleMembersRes{}, nil
}
