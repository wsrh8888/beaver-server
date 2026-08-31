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
	"time"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetFriendVerifiesListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetFriendVerifiesListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifiesListByIdsLogic {
	return &GetFriendVerifiesListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_friend_verifies_list_by_ids", ctx),
	}
}

func (l *GetFriendVerifiesListByIdsLogic) GetFriendVerifiesListByIds(in *friend_rpc.GetFriendVerifiesListByIdsReq) (*friend_rpc.GetFriendVerifiesListByIdsRes, error) {
	if len(in.VerifyIds) == 0 {
		return &friend_rpc.GetFriendVerifiesListByIdsRes{FriendVerifies: []*friend_rpc.FriendVerifyListById{}}, nil
	}

	// 查询指定ID列表中的好友验证信息
	var friendVerifies []friend_models.FriendVerifyModel
	query := l.svcCtx.DB.Where("verify_id IN (?)", in.VerifyIds)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&friendVerifies).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询好友验证信息失败", Data: map[string]any{"ids": in.VerifyIds, "since": in.Since, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{Text: "批量查询好友验证信息", Data: map[string]interface{}{"count": len(friendVerifies)}})

	// 转换为响应格式
	var friendVerifiesList []*friend_rpc.FriendVerifyListById
	for _, verify := range friendVerifies {
		friendVerifiesList = append(friendVerifiesList, &friend_rpc.FriendVerifyListById{
			VerifyId:   verify.VerifyID,
			SendUserId: verify.SendUserID,
			RevUserId:  verify.RevUserID,
			SendStatus: int32(verify.SendStatus),
			RevStatus:  int32(verify.RevStatus),
			Message:    verify.Message,
			Source:     verify.Source,
			Version:    verify.Version,
			CreatedAt:  time.Time(verify.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(verify.UpdatedAt).UnixMilli(),
		})
	}

	return &friend_rpc.GetFriendVerifiesListByIdsRes{FriendVerifies: friendVerifiesList}, nil
}
