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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *open_rpc.GetUserInfoReq) (*open_rpc.GetUserInfoRes, error) {
	var token open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", in.AccessToken).First(&token).Error; err != nil {
		return nil, errors.New("无效的访问令牌")
	}
	if time.Now().Unix() > token.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: token.UserID})
	if err != nil || userRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}
	info := userRes.UserInfo

	return &open_rpc.GetUserInfoRes{
		UserId:   info.UserId,
		NickName: info.NickName,
		Avatar:   info.Avatar,
		Phone:    info.Phone,
		Email:    info.Email,
	}, nil
}
