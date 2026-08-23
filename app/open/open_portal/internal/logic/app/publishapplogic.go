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

package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布应用（对标飞书版本发布）
func NewPublishAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishAppLogic {
	return &PublishAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishAppLogic) PublishApp(req *types.PublishAppReq) (resp *types.PublishAppRes, err error) {

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	if app.EnableRobot == 1 {
		if err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
			logx.Errorf("发布应用时 Robot 未就绪: app_id=%s err=%v", req.AppID, err)
			return nil, errors.New("发布失败：智能机器人未创建成功，请先开启 robot 能力后重试")
		}
	}

	// 更新应用状态为已发布
	if err := l.svcCtx.DB.Model(&app).Update("status", 1).Error; err != nil {
		logx.Errorf("发布应用失败: %v", err)
		return nil, errors.New("发布应用失败")
	}

	// 4. 记录版本信息（可选，后续可以创建 open_app_versions 表）
	logx.Infof("应用发布成功: app_id=%s, version=%s, user_id=%s", req.AppID, req.Version, req.UserID)

	return &types.PublishAppRes{
		Status: 1, // 已发布
	}, nil
}
