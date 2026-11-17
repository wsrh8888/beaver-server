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

package friendlogic

import (
	"context"
	"errors"

	friend_models "beaver/app/relation/relation_models/friend"
	"beaver/app/relation/relation_rpc/internal/svc"
	"beaver/app/relation/relation_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

const friendVerifyActionDelete int32 = 1 // 删除好友验证记录

type UpdateFriendVerifiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendVerifiesLogic {
	return &UpdateFriendVerifiesLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("update_friend_verifies", ctx)}
}

func (l *UpdateFriendVerifiesLogic) UpdateFriendVerifies(in *pb.UpdateFriendVerifiesReq) (*pb.UpdateFriendVerifiesRes, error) {
	if in.Action != friendVerifyActionDelete {
		return nil, errors.New("不支持的操作类型")
	}

	var ids []uint
	for _, vid := range in.VerifyIds {
		v, err := findFriendVerify(l.svcCtx.DB, vid)
		if err != nil {
			continue
		}
		ids = append(ids, v.ID)
	}
	if len(ids) == 0 {
		return &pb.UpdateFriendVerifiesRes{}, nil
	}

	if err := l.svcCtx.DB.Where("id IN ?", ids).Delete(&friend_models.FriendVerifyModel{}).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除好友验证失败", Data: map[string]any{"ids": ids, "err": err.Error()}})
		return nil, err
	}
	l.logger.Info(model.LogMsg{Text: "删除好友验证成功", Data: map[string]interface{}{"count": len(ids)}})
	return &pb.UpdateFriendVerifiesRes{AffectedCount: int64(len(ids))}, nil
}
