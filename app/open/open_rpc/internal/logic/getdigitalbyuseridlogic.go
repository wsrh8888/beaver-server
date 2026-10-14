/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
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

	var partner open_models.OpenDigital
	if err := l.svcCtx.DB.Where("digital_id = ?", in.DigitalUserId).First(&partner).Error; err != nil {
		return &open_rpc.GetDigitalByUserIDRes{Found: false}, nil
	}

	return &open_rpc.GetDigitalByUserIDRes{
		Found:         true,
		AppKey:        partner.AppKey,
		DigitalUserId: partner.DigitalID,
		Runtime:       partner.Runtime,
	}, nil
}
