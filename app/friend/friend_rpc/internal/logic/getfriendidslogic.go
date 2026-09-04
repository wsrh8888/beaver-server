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

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetFriendIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetFriendIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendIdsLogic {
	return &GetFriendIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_friend_ids", ctx),
	}
}

func (l *GetFriendIdsLogic) GetFriendIds(in *friend_rpc.GetFriendIdsRequest) (*friend_rpc.GetFriendIdsResponse, error) {
	var friends []friend_models.FriendModel
	err := l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = false", in.UserID, in.UserID).Find(&friends).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询好友ID列表失败", Data: map[string]any{"userId": in.UserID, "err": err.Error()}})
		return nil, err
	}

	var friendIds []string
	for _, friend := range friends {
		if friend.SendUserID == in.UserID {
			friendIds = append(friendIds, friend.RevUserID)
		} else {
			friendIds = append(friendIds, friend.SendUserID)
		}
	}

	return &friend_rpc.GetFriendIdsResponse{
		FriendIds: friendIds,
	}, nil
}
