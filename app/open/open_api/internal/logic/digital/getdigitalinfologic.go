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

package digital

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type GetDigitalInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDigitalInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalInfoLogic {
	return &GetDigitalInfoLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("get_digital_info", ctx)}
}

func (l *GetDigitalInfoLogic) GetDigitalInfo(authorization string) (resp *types.GetDigitalInfoRes, err error) {
	token, err := utils.ValidateDigitalAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	partner, err := utils.ResolveDigitalFromToken(l.svcCtx.DB, token.DigitalID)
	if err != nil {
		return nil, err
	}

	resp = &types.GetDigitalInfoRes{
		DigitalID: partner.DigitalID,
		AppKey:    partner.AppKey,
	}
	if userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: partner.DigitalID}); err == nil && userRes.UserInfo != nil {
		resp.DigitalName = userRes.UserInfo.NickName
		resp.Avatar = userRes.UserInfo.Avatar
	}
	return resp, nil
}
