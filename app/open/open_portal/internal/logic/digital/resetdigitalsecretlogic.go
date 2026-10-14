/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"

	"github.com/google/uuid"
)

type ResetDigitalSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResetDigitalSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetDigitalSecretLogic {
	return &ResetDigitalSecretLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("reset_digital_secret", ctx)}
}

func (l *ResetDigitalSecretLogic) ResetDigitalSecret(req *types.ResetDigitalSecretReq) (resp *types.ResetDigitalSecretRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}
	if partner.Runtime != open_models.DigitalRuntimeExternal {
		return nil, errors.New("仅外接数字伙伴可重置密钥")
	}

	newSecret := uuid.New().String() + uuid.New().String()
	if err := l.svcCtx.DB.Model(partner).Update("app_secret_hash", newSecret).Error; err != nil {
		return nil, errors.New("重置失败")
	}
	return &types.ResetDigitalSecretRes{AppSecret: newSecret}, nil
}
