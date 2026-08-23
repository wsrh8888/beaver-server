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
	"fmt"
	"strings"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIncomingWebhookLogic {
	return &CreateIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIncomingWebhookLogic) CreateIncomingWebhook(req *types.CreateIncomingWebhookReq) (resp *types.CreateIncomingWebhookRes, err error) {
	if req.AppID == "" || req.GroupID == "" {
		return nil, errors.New("appId 和 groupId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil || len(groupRes.Groups) == 0 {
		return nil, errors.New("群组不存在")
	}

	name := req.Name
	if name == "" {
		name = "通知机器人"
	}

	userRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user.UserCreateReq{
		NickName: name,
		UserType: int32(user_models.UserTypeBot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, errors.New("创建推送 Bot 用户失败")
	}

	rpcRes, err := l.svcCtx.OpenRpc.CreateBot(l.ctx, &open_rpc.CreateBotReq{
		GroupId: req.GroupID,
		BotId:   userRes.UserID,
	})
	if err != nil {
		return nil, errors.New("创建推送 Bot 失败")
	}

	if err := l.svcCtx.DB.Model(&open_models.OpenBotModel{}).Where("id = ?", rpcRes.Id).Updates(map[string]interface{}{
		"app_id": req.AppID,
		"name":   name,
	}).Error; err != nil {
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		return nil, errors.New("保存 Bot 元数据失败")
	}

	_, err = l.svcCtx.GroupRpc.AddGroupMember(l.ctx, &group_rpc.AddGroupMemberReq{
		GroupId:    req.GroupID,
		UserId:     userRes.UserID,
		OperatedBy: req.UserID,
	})
	if err != nil {
		_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: rpcRes.Id})
		logx.Errorf("推送 Bot 入群失败: group=%s bot=%s err=%v", req.GroupID, userRes.UserID, err)
		return nil, errors.New("推送 Bot 入群失败")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", rpcRes.Id).First(&bot).Error; err != nil {
		return nil, errors.New("查询 Bot 失败")
	}

	return &types.CreateIncomingWebhookRes{
		Webhook: toIncomingWebhookInfo(&bot, l.svcCtx.Config.Domain, true),
	}, nil
}

func toIncomingWebhookInfo(bot *open_models.OpenBotModel, apiBase string, withSecret bool) types.IncomingWebhookInfo {
	info := types.IncomingWebhookInfo{
		ID:         fmt.Sprintf("%d", bot.ID),
		Token:      bot.Token,
		AppID:      bot.AppID,
		GroupID:    bot.GroupID,
		BotID:      bot.BotID,
		Name:       bot.Name,
		WebhookURL: buildBotWebhookURL(apiBase, bot.Token),
		Status:     bot.Status,
		CreatedAt:  bot.CreatedAt.Unix(),
	}
	if withSecret && bot.Security.SignatureEnabled {
		info.Secret = bot.Security.SignatureSecret
	}
	return info
}

func buildBotWebhookURL(apiBase, token string) string {
	base := strings.TrimSuffix(apiBase, "/")
	return fmt.Sprintf("%s/api/open/bot_public/v1/send?token=%s", base, token)
}