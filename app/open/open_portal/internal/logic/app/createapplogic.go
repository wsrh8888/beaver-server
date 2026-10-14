/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package app

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/google/uuid"
)

type CreateAppLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建应用
func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{
		logger: beaverlog.New("create_app", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAppLogic) CreateApp(req *types.CreateAppReq) (resp *types.CreateAppRes, err error) {
	oauthID := fmt.Sprintf("oauth_%s", uuid.New().String()[:8])
	oauthSecret := uuid.New().String() + uuid.New().String()

	app := open_models.OpenOAuth{
		OAuthID:         oauthID,
		OAuthSecret:     oauthSecret,
		Name:            req.Name,
		Description:     req.Description,
		OwnerID:         req.UserID,
		Status:          1,
		Icon:            req.Icon,
		AccessTokenTTL:  7200,
		RefreshTokenTTL: 2592000,
	}

	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "创建应用失败", Data: map[string]interface{}{"err": err}})
		return nil, errors.New("创建应用失败")
	}

	l.logger.Info(model.LogMsg{Text: "应用创建成功", Data: map[string]interface{}{"oauth_id": oauthID, "user_id": req.UserID}})

	return &types.CreateAppRes{
		OAuthID:     oauthID,
		OAuthSecret: oauthSecret,
	}, nil
}
