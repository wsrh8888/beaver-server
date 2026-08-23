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

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const friendVerifyActionDelete int32 = 1 // 删除好友验证记录

type UpdateFriendVerifiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendVerifiesLogic {
	return &UpdateFriendVerifiesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendVerifiesLogic) UpdateFriendVerifies(in *friend_rpc.UpdateFriendVerifiesReq) (*friend_rpc.UpdateFriendVerifiesRes, error) {
	if in.Action != friendVerifyActionDelete {
		return nil, errors.New("不支持的操作类型")
	}

	var ids []uint
	for _, vid := range in.VerifyIds {
		v, err := findFriendVerify(l.svcCtx.DB, vid)
		if err != nil {
			continue
		}
		ids = append(ids, v.Id)
	}
	if len(ids) == 0 {
		return &friend_rpc.UpdateFriendVerifiesRes{}, nil
	}

	if err := l.svcCtx.DB.Where("id IN ?", ids).Delete(&friend_models.FriendVerifyModel{}).Error; err != nil {
		l.Errorf("删除好友验证失败: %v", err)
		return nil, err
	}
	return &friend_rpc.UpdateFriendVerifiesRes{AffectedCount: int64(len(ids))}, nil
}
