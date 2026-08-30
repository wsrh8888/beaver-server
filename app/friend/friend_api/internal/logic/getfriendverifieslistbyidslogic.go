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
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifiesListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取好友验证数据（通过ID）
func NewGetFriendVerifiesListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifiesListByIdsLogic {
	return &GetFriendVerifiesListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendVerifiesListByIdsLogic) GetFriendVerifiesListByIds(req *types.GetFriendVerifiesListByIdsReq) (resp *types.GetFriendVerifiesListByIdsRes, err error) {
	if len(req.VerifyIds) == 0 {
		return &types.GetFriendVerifiesListByIdsRes{
			FriendVerifies: []types.FriendVerifyById{},
		}, nil
	}

	// 查询指定ID列表中的好友验证信息
	var friendVerifies []friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Where("verify_id IN (?)", req.VerifyIds).Find(&friendVerifies).Error
	if err != nil {
		l.Errorf("查询好友验证信息失败: ids=%v, error=%v", req.VerifyIds, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友验证信息", len(friendVerifies))

	// 转换为响应格式
	var friendVerifiesList []types.FriendVerifyById
	for _, verify := range friendVerifies {
		friendVerifiesList = append(friendVerifiesList, types.FriendVerifyById{
			VerifyID:   verify.VerifyID,
			SendUserID: verify.SendUserID,
			RevUserID:  verify.RevUserID,
			SendStatus: int32(verify.SendStatus),
			RevStatus:  int32(verify.RevStatus),
			Message:    verify.Message,
			Source:     verify.Source,
			Version:    verify.Version,
			CreatedAt:  time.Time(verify.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(verify.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetFriendVerifiesListByIdsRes{
		FriendVerifies: friendVerifiesList,
	}, nil
}
