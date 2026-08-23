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

package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)


type GetH5AuthCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		ctx:    ctx,
		logger: logger.New("get_h5_auth_code"),
		svcCtx: svcCtx,
	}
}

func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (resp *types.GetH5AuthCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或未启用")
	}

	var oauthConfig open_models.OpenAppOAuth
	scope := ""
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scope = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scope = string(data)
	}

	const ttl = 5 * time.Minute
	authCode := util.NewV4().String()
	record := open_models.OpenOAuthCode{
		Code:      authCode,
		AppID:     req.AppID,
		UserID:    req.UserID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Scene:     "h5_sso",
	}
	if err := l.svcCtx.DB.Create(&record).Error; err != nil {
		logx.Errorf("生成 H5 authCode 失败: appId=%s, userId=%s, err=%v", req.AppID, req.UserID, err)
		return nil, errors.New("生成授权码失败")
	}

	logx.Infof("生成 H5 authCode 成功: appId=%s, userId=%s", req.AppID, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "H5授权码生成成功",
		Data: map[string]interface{}{
			"appId":  req.AppID,
			"userId": req.UserID,
		},
	})

	return &types.GetH5AuthCodeRes{
		AuthCode: authCode,
		ExpireIn: int64(ttl.Seconds()),
	}, nil
}
