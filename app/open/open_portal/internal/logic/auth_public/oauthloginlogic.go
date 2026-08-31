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

package auth_public

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/jwts"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type OAuthLoginLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOAuthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OAuthLoginLogic {
	return &OAuthLoginLogic{
		logger: beaverlog.New("oauth_login", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OAuthLoginLogic) OAuthLogin(req *types.OAuthLoginReq) (resp *types.OAuthLoginRes, err error) {
	if req.Code == "" {
		return nil, errors.New("授权码不能为空")
	}

	appID := l.svcCtx.Config.PortalOAuth.AppId
	if appID == "" {
		return nil, errors.New("门户 OAuth 未配置")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询门户 OAuth 应用失败", Data: map[string]interface{}{"app_id": appID, "err": err.Error()}})
		return nil, errors.New("门户 OAuth 应用不存在")
	}

	rpcResp, err := l.svcCtx.OpenRpc.ExchangeToken(l.ctx, &open_rpc.ExchangeTokenReq{
		AppId: appID,
		Code:  req.Code,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "code 换 token 失败", Data: map[string]interface{}{"app_id": appID, "err": err.Error()}})
		return nil, errors.New("授权码换取令牌失败")
	}
	if rpcResp.AccessToken == "" {
		return nil, errors.New("授权码换取令牌失败")
	}

	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ?", rpcResp.AccessToken).First(&tokenRecord).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询 OAuth token 失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("获取用户信息失败")
	}

	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: tokenRecord.UserID,
	})
	if err != nil || userRes.UserInfo == nil {
		l.logger.Error(model.LogMsg{Text: "查询用户信息失败", Data: map[string]interface{}{"user_id": tokenRecord.UserID, "err": err.Error()}})
		return nil, errors.New("获取用户信息失败")
	}

	secretKey := l.svcCtx.Config.Auth.AccessSecret
	expireHours := l.svcCtx.Config.Auth.AccessExpire / 3600
	if expireHours == 0 {
		expireHours = 12
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
	}, secretKey, int(expireHours))
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "生成 token 失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("服务内部异常")
	}

	expireAt := time.Now().Add(time.Duration(expireHours) * time.Hour).UnixMilli()
	l.logger.Info(model.LogMsg{Text: "OAuth 登录成功", Data: map[string]interface{}{"user_id": userRes.UserInfo.UserId, "nick_name": userRes.UserInfo.NickName}})

	return &types.OAuthLoginRes{
		Token:    token,
		UserID:   userRes.UserInfo.UserId,
		NickName: userRes.UserInfo.NickName,
		ExpireAt: expireAt,
	}, nil
}
