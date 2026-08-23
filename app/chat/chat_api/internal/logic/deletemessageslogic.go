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

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type DeleteMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 批量删除消息(仅自己不可见)
func NewDeleteMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessagesLogic {
	return &DeleteMessagesLogic{
		ctx:    ctx,
		logger: logger.New("delete_messages"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessagesLogic) DeleteMessages(req *types.DeleteMessagesReq) (resp *types.DeleteMessagesRes, err error) {
	// 1. 批量记录删除记录
	// 对标大厂：这里不改主表状态，而是插入标记表，让同步和历史记录接口过滤
	deleteRecords := make([]chat_models.ChatUserDelete, 0)
	for _, msgID := range req.MessageIDs {
		// 生成版本号 (对标大厂：每条操作都有独立版本，确保多端增量同步的原子性)
		version := l.svcCtx.VersionGen.GetNextVersion("chat_user_deletes", "user_id", req.UserID)
		deleteRecords = append(deleteRecords, chat_models.ChatUserDelete{
			UserID:    req.UserID,
			MessageID: msgID,
			Version:   version,
		})
	}

	// 3. 执行入库 (使用 Create 批量插入)
	err = l.svcCtx.DB.Create(&deleteRecords).Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("用户 %s 批量删除消息失败: %v", req.UserID, err)
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "批量删除消息成功",
		Data: map[string]interface{}{
			"userId": req.UserID,
			"count":  len(req.MessageIDs),
		},
	})

	return &types.DeleteMessagesRes{}, nil
}
