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
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/user"

	beaverlog "beaver/utils/beaverlog"
	"gorm.io/gorm"
)

type GetDigitalConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetDigitalConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDigitalConfigLogic {
	return &GetDigitalConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_digital_config", ctx),
	}
}

func (l *GetDigitalConfigLogic) GetDigitalConfig(req *types.GetDigitalConfigReq) (resp *types.GetDigitalConfigRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	partner, err := ensurePortalAppExternalPartner(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app)
	if err != nil {
		return nil, errors.New("获取数字伙伴配置失败")
	}

	config := types.DigitalConfigInfo{
		AppID:     req.AppID,
		DigitalID: partner.PartnerID,
		Status:    int(partner.Status),
	}
	if userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoReq{UserID: partner.PartnerID}); err == nil && userRes.UserInfo != nil {
		config.DigitalName = userRes.UserInfo.NickName
		config.Avatar = userRes.UserInfo.Avatar
	}

	return &types.GetDigitalConfigRes{Config: config}, nil
}

func ensurePortalAppExternalPartner(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) (*open_models.OpenDigitalPartner, error) {
	var partner open_models.OpenDigitalPartner
	err := db.Where(
		"app_id = ? AND runtime = ?",
		app.AppID, open_models.PartnerRuntimeExternal,
	).First(&partner).Error
	if err == nil && partner.PartnerID != "" {
		return &partner, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nickName := app.Name
	if nickName == "" {
		nickName = "数字伙伴"
	}

	createRes, err := userRpc.UserCreate(ctx, &user.UserCreateReq{
		NickName: nickName,
		UserType: int32(user_models.UserTypeDigital),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}

	partner = open_models.OpenDigitalPartner{
		PartnerID:  createRes.UserID,
		OwnerID:    app.OwnerID,
		Runtime:    open_models.PartnerRuntimeExternal,
		AppID:      app.AppID,
		Visibility: open_models.PartnerVisibilityPrivate,
		Status:     1,
	}
	if err := db.Create(&partner).Error; err != nil {
		return nil, err
	}
	return &partner, nil
}
