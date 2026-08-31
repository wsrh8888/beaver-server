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

package oauth_public

import (
	"context"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	util "beaver/utils/uuid"
)

type GenerateQrCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGenerateQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateQrCodeLogic {
	return &GenerateQrCodeLogic{
		ctx:    ctx,
		logger: beaverlog.New("generate_qrcode", ctx),
		svcCtx: svcCtx,
	}
}

func (l *GenerateQrCodeLogic) GenerateQrCode(req *types.GenerateQrCodeReq) (resp *types.GenerateQrCodeRes, err error) {
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "应用不存在",
			Data: map[string]any{"appId": req.AppID, "err": err.Error()},
		})
		return nil, fmt.Errorf("应用不存在")
	}
	if app.Status != 1 {
		l.logger.Warn(model.LogMsg{
			Text: "应用未启用",
			Data: map[string]any{"appId": req.AppID},
		})
		return nil, fmt.Errorf("应用未启用")
	}

	const expireIn int64 = 300
	sceneID := util.NewV4().String()
	expiresAt := time.Now().Add(time.Duration(expireIn) * time.Second)

	qrCode := open_models.OpenOAuthQrCode{
		SceneID:   sceneID,
		AppID:     req.AppID,
		Status:    0,
		ExpiresAt: expiresAt,
	}
	if err := l.svcCtx.DB.Create(&qrCode).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建扫码记录失败",
			Data: map[string]any{"appId": req.AppID, "err": err.Error()},
		})
		return nil, fmt.Errorf("服务内部异常")
	}

	l.logger.Info(model.LogMsg{
		Text: "生成扫码会话成功",
		Data: map[string]any{"sceneId": sceneID, "appId": req.AppID},
	})

	return &types.GenerateQrCodeRes{
		SceneID:  sceneID,
		ExpireIn: expireIn,
	}, nil
}
