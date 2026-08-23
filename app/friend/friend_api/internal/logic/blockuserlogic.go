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

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)


type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 拉黑用户
func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		ctx:    ctx,
		logger: logger.New("block_user"),
		svcCtx: svcCtx,
	}
}

func (l *BlockUserLogic) BlockUser(req *types.BlockUserReq) (resp *types.BlockUserRes, err error) {
	if req.UserID == req.BlockedUserID {
		return nil, errors.New("不能拉黑自己")
	}

	// 检查是否已拉黑
	var existing friend_models.FriendBlockModel
	result := l.svcCtx.DB.Where("user_id = ? AND blocked_user_id = ?", req.UserID, req.BlockedUserID).First(&existing)
	if result.Error == nil {
		return &types.BlockUserRes{}, nil // 已拉黑，幂等处理
	}

	block := friend_models.FriendBlockModel{
		BlockID:       uuid.New().String(),
		UserID:        req.UserID,
		BlockedUserID: req.BlockedUserID,
	}
	if err = l.svcCtx.DB.Create(&block).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("拉黑用户失败: userID=%s blockedUserID=%s err=%v", req.UserID, req.BlockedUserID, err)
		return nil, errors.New("操作失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "拉黑用户成功",
		Data: map[string]interface{}{
			"userId":        req.UserID,
			"blockedUserId": req.BlockedUserID,
		},
	})

	return &types.BlockUserRes{}, nil
}
