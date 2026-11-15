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

package auth_token

import (
	"context"
	"errors"

	"beaver/app/account/account_api/internal/svc"
	"beaver/app/account/account_api/internal/types"
	"beaver/commonconfig"
	"beaver/middleware/authgate"
	"beaver/middleware/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

// S2sTokenLogic 队间换票：client_credentials → 受众 core:api 的机器票。
//
// 这是本团队唯一的换票入口。原先由 datasync 承担（顺带挂了 74 条 RPC→HTTP 转发），
// 现在只留这一条小路由，能力本身回到各业务面自己的 /api 契约上。
type S2sTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewS2sTokenLogic 构造换票 logic。
func NewS2sTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *S2sTokenLogic {
	return &S2sTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// S2sToken 校验调用方凭据并签发机器票。
//
// fail-closed：注册表里没有该 client、Enabled 为假、口令不一致，一律拒绝；
// 口令比对走 jwt.AuthenticateClient 的常量时间实现，不做长度/前缀短路。
func (l *S2sTokenLogic) S2sToken(req *types.S2STokenReq) (resp *types.S2STokenRes, err error) {
	if req.GrantType != "" && req.GrantType != commonconfig.AuthClientCredentials {
		return nil, errors.New("不支持的授权方式")
	}

	scopes, ok := jwt.AuthenticateClient(l.svcCtx.Config.S2S.Clients, req.ClientID, req.ClientSecret)
	if !ok {
		l.Errorf("换票被拒：client_id=%s", req.ClientID)
		return nil, errors.New("invalid_client")
	}

	signed, err := jwt.Issue(
		l.svcCtx.Config.S2S.Secret,
		l.svcCtx.Config.Name,
		req.ClientID,
		authgate.AudienceCore,
		scopes,
		jwt.S2STTL,
	)
	if err != nil {
		l.Errorf("签发机器票失败：client_id=%s err=%v", req.ClientID, err)
		return nil, err
	}

	return &types.S2STokenRes{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(jwt.S2STTL.Seconds()),
	}, nil
}
