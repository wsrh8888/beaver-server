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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type DeleteNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 按事件ID删除单个通知
func NewDeleteNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotificationLogic {
	return &DeleteNotificationLogic{
		ctx:    ctx,
		logger: logger.New("delete_notification"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotificationLogic) DeleteNotification(req *types.DeleteNotificationReq) (resp *types.DeleteNotificationRes, err error) {
	userId := req.UserID
	eventId := req.EventID

	// 将指定用户的指定通知标记为已删除
	result := l.svcCtx.DB.Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND event_id = ?", userId, eventId).
		Update("is_deleted", true)

	if result.Error != nil {
		logx.WithContext(l.ctx).Errorf("删除通知失败: %v", result.Error)
		return nil, result.Error
	}

	resp = &types.DeleteNotificationRes{
		Success: result.RowsAffected > 0,
	}

	if result.RowsAffected > 0 {
		l.logger.Info(model.LogMsg{
			Text: "通知删除成功",
			Data: map[string]interface{}{
				"userId":  userId,
				"eventId": eventId,
			},
		})
	}

	return resp, nil
}
