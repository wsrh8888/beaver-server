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
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type CreateDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDigitalLogic {
	return &CreateDigitalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("create_digital", ctx),
	}
}

func (l *CreateDigitalLogic) CreateDigital(req *types.CreateDigitalReq) (resp *types.CreateDigitalRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	runtime := req.Runtime
	if runtime != open_models.DigitalRuntimeBuiltin &&
		runtime != open_models.DigitalRuntimeExternal &&
		runtime != open_models.DigitalRuntimeWorkflow {
		return nil, errors.New("runtime 无效")
	}
	visibility := req.Visibility
	if visibility == "" {
		visibility = open_models.DigitalVisibilityPrivate
	}

	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		NickName: req.Name,
		Abstract: req.Description,
		UserType: int32(user_models.UserTypeDigital),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, errors.New("创建数字伙伴 IM 用户失败")
	}
	if req.Avatar != "" {
		_, _ = l.svcCtx.UserRpc.UserUpdateDisplay(l.ctx, &user_rpc.UserUpdateDisplayReq{
			UserId: createRes.UserID,
			Avatar: req.Avatar,
		})
	}

	partner := open_models.OpenDigital{
		DigitalID:  createRes.UserID,
		OwnerID:    req.UserID,
		Runtime:    runtime,
		Visibility: visibility,
		Source:     open_models.DigitalSourceOpen,
	}

	res := &types.CreateDigitalRes{DigitalID: partner.DigitalID}
	if runtime == open_models.DigitalRuntimeExternal {
		appKey, appSecret := generateExternalCredentials()
		partner.AppKey = appKey
		partner.AppSecret = appSecret
		res.AppKey = appKey
		res.AppSecret = appSecret
	}

	if err := l.svcCtx.DB.Create(&partner).Error; err != nil {
		return nil, errors.New("创建数字伙伴失败")
	}
	return res, nil
}
