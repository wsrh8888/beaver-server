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

type GetOpenOAuthListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOpenOAuthListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOpenOAuthListLogic {
	return &GetOpenOAuthListLogic{logger: beaverlog.New("get_open_oauth_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOpenOAuthListLogic) GetOpenOAuthList(req *types.GetOpenOAuthListReq) (resp *types.GetOpenOAuthListRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.ListOpenOAuth(l.ctx, &open_rpc.ListOpenOAuthReq{
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
		Keyword:     req.Keyword,
		OwnerUserId: req.OwnerUserID,
		OauthId:     req.OAuthID,
		Status:      int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "获取 OAuth 应用列表失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, err
	}

	users := l.loadUsers(rpcRes.List)

	list := make([]types.OpenOAuthInfo, 0, len(rpcRes.List))
	for _, app := range rpcRes.List {
		ownerName := ""
		if u, ok := users[app.OwnerUserId]; ok && u != nil {
			ownerName = u.NickName
		}
		list = append(list, types.OpenOAuthInfo{
			OAuthID:       app.OauthId,
			Name:          app.Name,
			Description:   app.Description,
			Icon:          app.Icon,
			OwnerUserID:   app.OwnerUserId,
			OwnerUserName: ownerName,
			Status:        int(app.Status),
			CreatedAt:     app.CreatedAt,
		})
	}

	return &types.GetOpenOAuthListRes{Total: rpcRes.Total, List: list}, nil
}

// loadUsers 批量取所有者昵称，失败不影响列表本身。
func (l *GetOpenOAuthListLogic) loadUsers(apps []*open_rpc.OpenOAuthItem) map[string]*user_rpc.UserInfo {
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(apps))
	for _, app := range apps {
		if app.OwnerUserId == "" {
			continue
		}
		if _, ok := seen[app.OwnerUserId]; ok {
			continue
		}
		seen[app.OwnerUserId] = struct{}{}
		userIDs = append(userIDs, app.OwnerUserId)
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
