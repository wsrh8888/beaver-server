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

	"github.com/zeromicro/go-zero/core/logx"
)

var _ = logx.WithContext

type SetDigitalAdminsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSetDigitalAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDigitalAdminsLogic {
	return &SetDigitalAdminsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SetDigitalAdmins", ctx),
	}
}

func (l *SetDigitalAdminsLogic) SetDigitalAdmins(req *types.SetDigitalAdminsReq) (resp *types.SetDigitalAdminsRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}

	// diff 增删：仅操作 role=admin 行，不影响 role=viewer，未变更的不动。
	tx := l.svcCtx.DB.Begin()
	if err := syncDigitalMembers(tx, req.DigitalID, open_models.DigitalMemberRoleAdmin, req.UserIDs, req.UserID); err != nil {
		tx.Rollback()
		return nil, errors.New("更新管理员失败")
	}
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("更新管理员失败")
	}
	return &types.SetDigitalAdminsRes{}, nil
}
