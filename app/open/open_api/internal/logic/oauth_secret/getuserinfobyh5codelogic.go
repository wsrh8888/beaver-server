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

package oauth_secret

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByH5CodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoByH5CodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByH5CodeLogic {
	return &GetUserInfoByH5CodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByH5CodeLogic) GetUserInfoByH5Code(req *types.GetUserInfoByH5CodeReq) (resp *types.GetUserInfoByH5CodeRes, err error) {
	if _, err := verifyApp(l.svcCtx.DB, req.AppID, req.AppSecret); err != nil {
		return nil, err
	}

	oauthCode, err := findOAuthCode(l.svcCtx.DB, req.AppID, req.AuthCode)
	if err != nil {
		return nil, err
	}
	if oauthCode.Scene != "h5_sso" {
		return nil, errors.New("授权码场景无效")
	}

	if err := l.svcCtx.DB.Model(oauthCode).Update("used", true).Error; err != nil {
		logx.Errorf("标记 authCode 已使用失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	userInfoRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: oauthCode.UserID,
	})
	if err != nil {
		logx.Errorf("查询用户信息失败: %v", err)
		return nil, errors.New("获取用户信息失败")
	}
	if userInfoRes.UserInfo == nil {
		return nil, errors.New("用户不存在")
	}

	return &types.GetUserInfoByH5CodeRes{
		UserID:   userInfoRes.UserInfo.UserId,
		NickName: userInfoRes.UserInfo.NickName,
		Avatar:   userInfoRes.UserInfo.Avatar,
		Phone:    "",
		Email:    userInfoRes.UserInfo.Email,
	}, nil
}
