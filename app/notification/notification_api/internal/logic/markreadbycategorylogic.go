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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类标记所有通知为已读
func NewMarkReadByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByCategoryLogic {
	return &MarkReadByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByCategoryLogic) MarkReadByCategory(req *types.MarkReadByCategoryReq) (resp *types.MarkReadByCategoryRes, err error) {
	userId := req.UserID
	category := req.Category

	now := time.Now()

	// 生成游标版本号（按用户+分类分区）
	cursorVersion := l.svcCtx.VersionGen.GetNextVersion(notification_models.VersionScopeCursorPerUser, "user_id", userId)
	if cursorVersion == -1 {
		l.Logger.Errorf("生成游标版本号失败")
		return nil, errors.New("生成版本号失败")
	}

	// 更新或创建游标记录：优先更新，不存在则创建
	result := l.svcCtx.DB.Model(&notification_models.NotificationRead{}).
		Where("user_id = ? AND category = ?", userId, category).
		Updates(map[string]interface{}{
			"version":      cursorVersion,
			"last_read_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		l.Logger.Errorf("更新游标失败: %v", result.Error)
		return nil, result.Error
	}

	// 如果没有更新到记录，说明记录不存在，需要创建
	if result.RowsAffected == 0 {
		cursor := &notification_models.NotificationRead{
			UserID:     userId,
			Category:   category,
			Version:    cursorVersion,
			LastReadAt: &now,
		}

		err = l.svcCtx.DB.Create(cursor).Error
		if err != nil {
			l.Logger.Errorf("创建游标失败: %v", err)
			return nil, err
		}
	}

	// 异步通过 WS 通知用户其他客户端
	go func(etcdAddr string, userId string, category string, cursorVersion int64) {
		// 构建表更新数据
		var tableUpdates []map[string]interface{}

		// 通知已读游标表更新
		cursorUpdates := map[string]interface{}{
			"table":  "notification_read_cursor",
			"userId": userId,
			"data": []map[string]interface{}{
				{
					"version":  cursorVersion,
					"category": category,
				},
			},
		}
		tableUpdates = append(tableUpdates, cursorUpdates)

		// 通知给自己（用户ID作为接收者，空字符串作为发送者表示系统操作）
		payload := map[string]interface{}{
			"command":  wsCommandConst.NOTIFICATION,
			"type":     wsTypeConst.NotificationMarkReadReceive,
			"senderId": "",
			"targetId": userId,
			"body": map[string]interface{}{
				"tableUpdates": tableUpdates,
			},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
	}(l.svcCtx.Config.Etcd, userId, category, cursorVersion)

	resp = &types.MarkReadByCategoryRes{}

	return resp, nil
}
