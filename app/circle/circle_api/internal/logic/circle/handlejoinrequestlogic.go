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

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleJoinRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleJoinRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleJoinRequestLogic {
	return &HandleJoinRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleJoinRequestLogic) HandleJoinRequest(req *types.HandleJoinRequestReq) (resp *types.HandleJoinRequestRes, err error) {
	// 权限校验
	var joinReq circle_models.CircleJoinRequestModel
	if err = l.svcCtx.DB.Where("id = ? AND status = 0", req.RequestID).First(&joinReq).Error; err != nil {
		return nil, fmt.Errorf("申请记录不存在或已处理")
	}

	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", joinReq.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可处理申请")
	}

	// 更新申请状态
	if err = l.svcCtx.DB.Model(&joinReq).Update("status", req.Status).Error; err != nil {
		return nil, fmt.Errorf("处理申请失败: %v", err)
	}

	// 通过则加入成员
	if req.Status == 1 {
		member := circle_models.CircleMemberModel{
			CircleID: joinReq.CircleID,
			UserID:   joinReq.UserID,
			Role:     3,
		}
		l.svcCtx.DB.Create(&member)

		circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", joinReq.CircleID)
		l.svcCtx.DB.Model(&circle_models.CircleModel{}).
			Where("circle_id = ?", joinReq.CircleID).
			Update("version", circleVersion)

		conversationID := fmt.Sprintf("circle_%s", joinReq.CircleID)

		go func() {
			ctx := context.Background()
			l.svcCtx.ChatRpc.InitializeConversation(ctx, &chat_rpc.InitializeConversationReq{
				ConversationId: conversationID,
				Type:           3,
				UserIds:        []string{joinReq.UserID},
			})

			payload := map[string]interface{}{
				"command":  wsCommandConst.CIRCLE_OPERATION,
				"type":     wsTypeConst.CircleReceive,
				"senderId": req.UserID,
				"targetId": joinReq.UserID,
				"body": map[string]interface{}{
					"tables": []map[string]interface{}{
						{
							"table": "circles",
							"data": []map[string]interface{}{
								{
									"version":  circleVersion,
									"circleId": joinReq.CircleID,
								},
							},
						},
					},
				},
				"conversationId": conversationID,
			}
			if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
				logx.Errorf("推送圈子资料同步失败: circleID=%s, err=%v", joinReq.CircleID, err)
			}
		}()
	}

	return &types.HandleJoinRequestRes{}, nil
}
