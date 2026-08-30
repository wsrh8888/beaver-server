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

package event

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEventSubscriptionLogic {
	return &DeleteEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEventSubscriptionLogic) DeleteEventSubscription(req *types.DeleteEventSubscriptionReq) (resp *types.DeleteEventSubscriptionRes, err error) {
	subID, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil || subID == 0 {
		return nil, errors.New("订阅 ID 无效")
	}

	var sub open_models.OpenAppEventSubscription
	if err := l.svcCtx.DB.Where("id = ?", subID).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订阅不存在")
		}
		return nil, errors.New("查询订阅失败")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", sub.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	result := l.svcCtx.DB.Where("id = ?", subID).Delete(&open_models.OpenAppEventSubscription{})
	if result.Error != nil {
		return nil, errors.New("删除订阅失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("订阅不存在")
	}

	return &types.DeleteEventSubscriptionRes{}, nil
}
