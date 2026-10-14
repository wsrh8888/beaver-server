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
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type UpdateOpenOAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateOpenOAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOpenOAuthLogic {
	return &UpdateOpenOAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateOpenOAuthLogic", ctx),
	}
}

func (l *UpdateOpenOAuthLogic) UpdateOpenOAuth(in *open_rpc.UpdateOpenOAuthReq) (*open_rpc.UpdateOpenOAuthRes, error) {
	if len(in.OauthIds) == 0 {
		return nil, errors.New("oauth_ids 不能为空")
	}
	if in.Status != 0 && in.Status != 1 {
		return nil, errors.New("status 无效")
	}

	result := l.svcCtx.DB.Model(&open_models.OpenOAuth{}).
		Where("oauth_id IN ?", in.OauthIds).
		Update("status", in.Status)
	if result.Error != nil {
		return nil, result.Error
	}
	return &open_rpc.UpdateOpenOAuthRes{AffectedCount: result.RowsAffected}, nil
}
