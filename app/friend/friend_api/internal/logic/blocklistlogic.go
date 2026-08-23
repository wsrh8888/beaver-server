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
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取黑名单列表
func NewBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockListLogic {
	return &BlockListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlockListLogic) BlockList(req *types.BlockListReq) (resp *types.BlockListRes, err error) {
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	var blocks []friend_models.FriendBlockModel
	var count int64
	query := l.svcCtx.DB.Model(&friend_models.FriendBlockModel{}).Where("user_id = ?", req.UserID)
	if err = query.Count(&count).Error; err != nil {
		return nil, errors.New("查询失败")
	}
	if err = query.Offset(offset).Limit(limit).Find(&blocks).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	list := make([]types.BlockUserInfo, 0, len(blocks))
	for _, b := range blocks {
		info := types.BlockUserInfo{
			BlockID: b.BlockID,
			UserID:  b.BlockedUserID,
		}
		userResp, uErr := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: b.BlockedUserID})
		if uErr == nil && userResp.UserInfo != nil {
			info.NickName = userResp.UserInfo.NickName
			info.Avatar = userResp.UserInfo.Avatar
		}
		list = append(list, info)
	}

	return &types.BlockListRes{List: list, Count: count}, nil
}
