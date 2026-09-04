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

type GetFriendVerifyVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetFriendVerifyVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyVersionsLogic {
	return &GetFriendVerifyVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_friend_verify_versions", ctx),
	}
}

func (l *GetFriendVerifyVersionsLogic) GetFriendVerifyVersions(in *friend_rpc.GetFriendVerifyVersionsReq) (*friend_rpc.GetFriendVerifyVersionsRes, error) {
	// 查询用户相关的所有好友验证记录（作为发送者或接收者）
	var friendVerifies []friend_models.FriendVerifyModel
	query := l.svcCtx.DB.Where("send_user_id = ? OR rev_user_id = ?", in.UserId, in.UserId)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&friendVerifies).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询用户好友验证版本信息失败", Data: map[string]any{"userId": in.UserId, "since": in.Since, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "查询用户好友验证版本信息", Data: map[string]interface{}{"userId": in.UserId, "count": len(friendVerifies)}})

	// 转换为响应格式
	var friendVerifyVersions []*friend_rpc.GetFriendVerifyVersionsRes_FriendVerifyVersion
	for _, verify := range friendVerifies {
		friendVerifyVersions = append(friendVerifyVersions, &friend_rpc.GetFriendVerifyVersionsRes_FriendVerifyVersion{
			VerifyId: verify.VerifyID,
			Version:  verify.Version,
		})
	}

	return &friend_rpc.GetFriendVerifyVersionsRes{FriendVerifyVersions: friendVerifyVersions}, nil
}
