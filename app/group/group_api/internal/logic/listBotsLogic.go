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
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBotsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群内所有机器人列表
func NewListBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBotsLogic {
	return &ListBotsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBotsLogic) ListBots(req *types.ListBotsReq) (resp *types.ListBotsRes, err error) {
	var member group_models.GroupMemberModel
	if err = l.svcCtx.DB.Take(&member, "group_id = ? AND user_id = ?", req.GroupID, req.UserID).Error; err != nil {
		return nil, errors.New("不是群成员")
	}
	if member.Role != 1 && member.Role != 2 {
		return nil, errors.New("无权限，仅群主或管理员可查看机器人")
	}

	// 查询群内机器人
	var bots []group_models.GroupBotModel
	if err = l.svcCtx.DB.Where("group_id = ?", req.GroupID).
		Order("id DESC").Find(&bots).Error; err != nil {
		return nil, err
	}

	// 批量获取用户信息（通过 user_rpc）
	botIDs := make([]string, 0, len(bots))
	for _, b := range bots {
		botIDs = append(botIDs, b.BotID)
	}

	userRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: botIDs,
	})
	if err != nil {
		return nil, err
	}

	items := make([]types.ListBotsItem, 0, len(bots))
	for _, b := range bots {
		userInfo := userRes.UserInfo[b.BotID]
		if userInfo == nil {
			continue
		}

		items = append(items, types.ListBotsItem{
			BotID:       b.BotID,
			Name:        userInfo.NickName,
			Description: userInfo.Abstract,
			Avatar:      userInfo.Avatar,
			Type:        b.Type,
			Status:      b.Status,
			CreatedAt:   time.Time(b.CreatedAt).Unix(),
		})
	}

	return &types.ListBotsRes{List: items}, nil
}
