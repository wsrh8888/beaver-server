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
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type JoinCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinCircleLogic {
	return &JoinCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinCircleLogic) JoinCircle(req *types.JoinCircleReq) (resp *types.JoinCircleRes, err error) {
	var invite *circle_models.CircleInviteModel
	circleID := req.CircleID
	if req.InviteCode != "" {
		var row circle_models.CircleInviteModel
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
		if circleID != "" && circleID != invite.CircleID {
			return nil, fmt.Errorf("邀请与圈子不匹配")
		}
		circleID = invite.CircleID
	}
	if circleID == "" {
		return nil, fmt.Errorf("圈子ID不能为空")
	}

	var circle circle_models.CircleModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", circleID).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	var existing circle_models.CircleMemberModel
	if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", circleID, req.UserID).First(&existing).Error == nil {
		return &types.JoinCircleRes{Status: 1, CircleID: circleID}, nil
	}

	if circle.JoinType == 1 {
		joinReq := circle_models.CircleJoinRequestModel{
			CircleID: circleID,
			UserID:   req.UserID,
			Status:   0,
			Reason:   req.Reason,
		}
		if err = l.svcCtx.DB.Create(&joinReq).Error; err != nil {
			return nil, fmt.Errorf("提交申请失败: %v", err)
		}
		l.bumpInviteUse(invite)
		return &types.JoinCircleRes{Status: 0, CircleID: circleID}, nil
	}

	member := circle_models.CircleMemberModel{
		CircleID: circleID,
		UserID:   req.UserID,
		Role:     3,
	}
	if err = l.svcCtx.DB.Create(&member).Error; err != nil {
		return nil, fmt.Errorf("加入圈子失败: %v", err)
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", circleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", circleID).
		Update("version", circleVersion)

	l.bumpInviteUse(invite)

	conversationID := fmt.Sprintf("circle_%s", circleID)
	go func() {
		ctx := context.Background()
		l.svcCtx.ChatRpc.InitializeConversation(ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        []string{req.UserID},
		})
		payload := map[string]interface{}{
			"command":  wsCommandConst.CIRCLE_OPERATION,
			"type":     wsTypeConst.CircleReceive,
			"senderId": req.UserID,
			"targetId": req.UserID,
			"body": map[string]interface{}{
				"tables": []map[string]interface{}{
					{
						"table": "circles",
						"data": []map[string]interface{}{
							{
								"version":  circleVersion,
								"circleId": circleID,
							},
						},
					},
				},
			},
			"conversationId": conversationID,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logx.Errorf("推送圈子资料同步失败: circleID=%s, err=%v", circleID, err)
		}
	}()

	return &types.JoinCircleRes{Status: 1, CircleID: circleID}, nil
}

func (l *JoinCircleLogic) bumpInviteUse(invite *circle_models.CircleInviteModel) {
	if invite == nil {
		return
	}
	updates := map[string]interface{}{
		"used_count": gorm.Expr("used_count + 1"),
	}
	if invite.MaxUses > 0 && invite.UsedCount+1 >= invite.MaxUses {
		updates["status"] = 3
	}
	_ = l.svcCtx.DB.Model(&circle_models.CircleInviteModel{}).Where("id = ?", invite.Id).Updates(updates).Error
}
