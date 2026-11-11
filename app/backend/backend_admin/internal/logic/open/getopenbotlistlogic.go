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

package open

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetOpenBotListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOpenBotListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOpenBotListLogic {
	return &GetOpenBotListLogic{logger: beaverlog.New("get_open_bot_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOpenBotListLogic) GetOpenBotList(req *types.GetOpenBotListReq) (resp *types.GetOpenBotListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListBots(l.ctx, &open_rpc.ListBotsReq{
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
		Keyword:     req.Keyword,
		Status:      int32(req.Status),
		OwnerUserId: req.OwnerUserID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "获取推送机器人列表失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, err
	}

	users := l.loadUsers(rpcRes.List)

	list := make([]types.OpenBotInfo, 0, len(rpcRes.List))
	for _, bot := range rpcRes.List {
		botName := ""
		if u, ok := users[bot.BotId]; ok && u != nil {
			botName = u.NickName
		}
		ownerName := ""
		if u, ok := users[bot.OwnerId]; ok && u != nil {
			ownerName = u.NickName
		}
		list = append(list, types.OpenBotInfo{
			BotID:            bot.BotId,
			BotName:          botName,
			OwnerID:          bot.OwnerId,
			OwnerUserName:    ownerName,
			Visibility:       bot.Visibility,
			Status:           int(bot.Status),
			SignatureEnabled: int(bot.SignatureEnabled),
			CreatedAt:        bot.CreatedAt,
		})
	}

	return &types.GetOpenBotListRes{Total: rpcRes.Total, List: list}, nil
}

// loadUsers 批量取 Bot 与创建者昵称，失败不影响列表本身。
func (l *GetOpenBotListLogic) loadUsers(bots []*open_rpc.BotItem) map[string]*user_rpc.UserInfo {
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(bots)*2)
	add := func(id string) {
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		userIDs = append(userIDs, id)
	}
	for _, bot := range bots {
		add(bot.BotId)
		add(bot.OwnerId)
	}
	if len(userIDs) == 0 {
		return map[string]*user_rpc.UserInfo{}
	}
	res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})
	if err != nil || res == nil {
		return map[string]*user_rpc.UserInfo{}
	}
	return res.UserInfo
}
