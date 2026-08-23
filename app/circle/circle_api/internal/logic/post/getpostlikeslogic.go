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

package post

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLikesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLikesLogic {
	return &GetPostLikesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostLikesLogic) GetPostLikes(req *types.GetPostLikesReq) (resp *types.GetPostLikesRes, err error) {
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	offset := (page - 1) * limit

	var totalCount int64
	if err = l.svcCtx.DB.Model(&circle_models.CircleLikeModel{}).
		Where("post_id = ?", req.PostID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	var likes []circle_models.CircleLikeModel
	if err = l.svcCtx.DB.Where("post_id = ?", req.PostID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	userIds := make([]string, 0, len(likes))
	for _, like := range likes {
		userIds = append(userIds, like.UserID)
	}

	userInfoMap := make(map[string]*user_rpc.UserInfo)
	if len(userIds) > 0 {
		userResp, rpcErr := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if rpcErr == nil && userResp != nil {
			userInfoMap = userResp.UserInfo
		}
	}

	likeInfos := make([]types.GetPostLikesInfo, 0, len(likes))
	for _, like := range likes {
		userName, avatar := "", ""
		if info := userInfoMap[like.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
		likeInfos = append(likeInfos, types.GetPostLikesInfo{
			UserID:    like.UserID,
			UserName:  userName,
			Avatar:    avatar,
			CreatedAt: like.CreatedAt.String(),
		})
	}

	return &types.GetPostLikesRes{
		Count: totalCount,
		List:  likeInfos,
	}, nil
}
