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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type GetDigitalByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDigitalByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalByUserIDLogic {
	return &GetDigitalByUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_digital_by_user_id", ctx),
	}
}

func (l *GetDigitalByUserIDLogic) GetDigitalByUserID(in *open_rpc.GetDigitalByUserIDReq) (*open_rpc.GetDigitalByUserIDRes, error) {
	if in.DigitalUserId == "" {
		return &open_rpc.GetDigitalByUserIDRes{Found: false}, nil
	}

	var partner open_models.OpenDigitalPartner
	if err := l.svcCtx.DB.Where("partner_id = ? AND status = 1", in.DigitalUserId).First(&partner).Error; err != nil {
		return &open_rpc.GetDigitalByUserIDRes{Found: false}, nil
	}

	return &open_rpc.GetDigitalByUserIDRes{
		Found:          true,
		AppId:          partner.AppID,
		DigitalUserId:  partner.PartnerID,
	}, nil
}
