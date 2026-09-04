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

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("user-info", ctx),
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoRes, err error) {
	l.logger.Info(model.LogMsg{
		Text: "获取用户基础信息",
		Data: map[string]any{"userId": req.UserID},
	})

	// 直接从数据库查询，避免RPC调用自身服务
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "user_id = ?", req.UserID).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询用户失败",
			Data: map[string]any{"userId": req.UserID, "err": err.Error()},
		})
		return nil, err
	}

	resp = &types.UserInfoRes{
		UserID:   user.UserID,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Abstract: user.Abstract,
		Phone:    user.Phone,
		Email:    user.Email,
		Gender:   user.Gender,
		Version:  user.Version,
	}

	return resp, nil
}
