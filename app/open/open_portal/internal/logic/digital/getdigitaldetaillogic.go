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
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type GetDigitalDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDigitalDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalDetailLogic {
	return &GetDigitalDetailLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("get_digital_detail", ctx)}
}

func (l *GetDigitalDetailLogic) GetDigitalDetail(req *types.GetDigitalDetailReq) (resp *types.GetDigitalDetailRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}
	info := fillDigitalInfo(l.ctx, l.svcCtx, partner)
	return &types.GetDigitalDetailRes{Digital: info}, nil
}
