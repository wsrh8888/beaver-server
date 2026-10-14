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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type UpdateDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDigitalLogic {
	return &UpdateDigitalLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("update_digital", ctx)}
}

func (l *UpdateDigitalLogic) UpdateDigital(req *types.UpdateDigitalReq) (resp *types.UpdateDigitalRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" || req.Avatar != "" {
		_, err = l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId:   partner.DigitalID,
			NickName: req.Name,
			Avatar:   req.Avatar,
		})
		if err != nil {
			return nil, errors.New("更新展示信息失败")
		}
	}
	if req.Description != "" {
		abstract := req.Description
		_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, &user_rpc.UpdateUsersReq{
			UserIds:       []string{partner.DigitalID},
			PatchAbstract: &abstract,
		})
		if err != nil {
			return nil, errors.New("更新简介失败")
		}
	}

	if req.Visibility != "" {
		if err := l.svcCtx.DB.Model(partner).Update("visibility", req.Visibility).Error; err != nil {
			return nil, errors.New("更新失败")
		}
	}
	return &types.UpdateDigitalRes{}, nil
}
