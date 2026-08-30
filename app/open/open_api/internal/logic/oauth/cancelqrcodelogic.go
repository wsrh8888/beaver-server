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
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CancelQrCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCancelQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelQrCodeLogic {
	return &CancelQrCodeLogic{
		ctx:    ctx,
		logger: beaverlog.New("cancel_qrcode", ctx),
		svcCtx: svcCtx,
	}
}

func (l *CancelQrCodeLogic) CancelQrCode(req *types.CancelQrCodeReq) (resp *types.CancelQrCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.SceneID == "" {
		return nil, errors.New("sceneId 不能为空")
	}

	if err := l.svcCtx.OAuth.Cancel(req.SceneID, req.UserID); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "扫码取消失败",
			Data: map[string]any{"sceneId": req.SceneID, "userId": req.UserID, "err": err.Error()},
		})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "OAuth扫码取消成功",
		Data: map[string]any{
			"sceneId": req.SceneID,
			"userId":  req.UserID,
		},
	})

	return &types.CancelQrCodeRes{Success: true}, nil
}
