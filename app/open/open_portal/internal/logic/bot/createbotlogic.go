/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package bot

import (
	"context"
	"errors"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type CreateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBotLogic {
	return &CreateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("portal_create_bot", ctx),
	}
}

func (l *CreateBotLogic) CreateBot(req *types.CreateBotReq) (*types.CreateBotRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = open_models.BotVisibilityPrivate
	}

	if req.GroupID != "" {
		groupRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, &group_rpc.GetGroupsListByIdsReq{
			GroupIDs: []string{req.GroupID},
		})
		if err != nil || len(groupRes.Groups) == 0 {
			return nil, errors.New("群组不存在")
		}
	}

	userRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		NickName: req.Name,
		Abstract: req.Description,
		UserType: int32(user_models.UserTypeBot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, errors.New("创建机器人 IM 用户失败")
	}
	if req.Avatar != "" {
		_, _ = l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId: userRes.UserID,
			Avatar: req.Avatar,
		})
	}

	rpcRes, err := l.svcCtx.OpenRpc.CreateBot(l.ctx, &open_rpc.CreateBotReq{
		BotId:      userRes.UserID,
		OwnerId:    req.UserID,
		Visibility: visibility,
	})
	if err != nil {
		return nil, errors.New("创建机器人失败")
	}

	if req.GroupID != "" {
		_, err = l.svcCtx.GroupRpc.AddGroupMember(l.ctx, &group_rpc.AddGroupMemberReq{
			GroupId:    req.GroupID,
			UserId:     userRes.UserID,
			OperatedBy: req.UserID,
		})
		if err != nil {
			_, _ = l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{BotId: userRes.UserID})
			return nil, errors.New("机器人入群失败")
		}
	}

	return &types.CreateBotRes{
		BotID:      userRes.UserID,
		WebhookURL: buildBotWebhookURL(l.svcCtx.Config.Domain, rpcRes.Token),
		Token:      rpcRes.Token,
		Secret:     rpcRes.SignatureSecret,
	}, nil
}
