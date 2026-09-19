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

package bot

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/user"
	beaverlog "beaver/utils/beaverlog"
)

type ListIncomingWebhooksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListIncomingWebhooksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIncomingWebhooksLogic {
	return &ListIncomingWebhooksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("list_incoming_webhooks", ctx),
	}
}

func (l *ListIncomingWebhooksLogic) ListIncomingWebhooks(req *types.ListIncomingWebhooksReq) (resp *types.ListIncomingWebhooksRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenBotModel{}).Where("app_id = ?", req.AppID)
	if req.GroupID != "" {
		query = query.Where("group_id = ?", req.GroupID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	var bots []open_models.OpenBotModel
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bots).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	// 显示名不再存储在机器人表，批量取 IM 用户资料（BotID 即 IM 用户ID）
	names := make(map[string]string, len(bots))
	botIDs := make([]string, 0, len(bots))
	for i := range bots {
		if bots[i].BotID != "" {
			botIDs = append(botIDs, bots[i].BotID)
		}
	}
	if len(botIDs) > 0 {
		if userRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user.UserListInfoReq{UserIdList: botIDs}); err == nil && userRes.UserInfo != nil {
			for id, info := range userRes.UserInfo {
				if info != nil {
					names[id] = info.NickName
				}
			}
		}
	}

	list := make([]types.IncomingWebhookInfo, 0, len(bots))
	for i := range bots {
		list = append(list, toIncomingWebhookInfo(&bots[i], l.svcCtx.Config.Domain, false, names[bots[i].BotID]))
	}

	return &types.ListIncomingWebhooksRes{
		Total: total,
		List:  list,
	}, nil
}