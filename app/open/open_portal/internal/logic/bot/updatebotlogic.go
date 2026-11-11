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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type UpdateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotLogic {
	return &UpdateBotLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("portal_update_bot", ctx)}
}

func (l *UpdateBotLogic) UpdateBot(req *types.UpdateBotReq) (*types.UpdateBotRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedBot(l.svcCtx.DB, req.UserID, req.BotID)
	if err != nil {
		return nil, err
	}

	// 名称/头像/描述都存在 IM 用户资料里（BotID 即 UserModel.UserID），不复制到 open_bot：
	// 复制一份就要双写两份、还得手动同步，绑定关系用 ID 表达即可。
	// 用 UpdateUsers 而不是 UserUpdateDisplay，是因为后者不支持改简介，且这里要能区分「不传」和「传空串清空」。
	if req.Name != nil || req.Avatar != nil || req.Description != nil {
		if _, err := l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
			UserIds:       []string{row.BotID},
			PatchNickName: req.Name,
			PatchAvatar:   req.Avatar,
			PatchAbstract: req.Description,
		}); err != nil {
			return nil, errors.New("更新机器人展示信息失败")
		}
	}

	if req.Visibility != "" {
		if _, err := l.svcCtx.OpenRpc.UpdateBot(l.ctx, &open_rpc.UpdateBotReq{
			BotId:      row.BotID,
			Visibility: req.Visibility,
		}); err != nil {
			return nil, errors.New("更新机器人失败")
		}
	}

	return &types.UpdateBotRes{}, nil
}
