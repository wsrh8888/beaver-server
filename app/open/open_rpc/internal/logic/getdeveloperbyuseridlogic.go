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

package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDeveloperByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeveloperByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperByUserIDLogic {
	return &GetDeveloperByUserIDLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetDeveloperByUserIDLogic) GetDeveloperByUserID(in *open_rpc.GetDeveloperByUserIDReq) (*open_rpc.GetDeveloperByUserIDRes, error) {
	var dev open_models.OpenDeveloper
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &open_rpc.GetDeveloperByUserIDRes{Found: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return &open_rpc.GetDeveloperByUserIDRes{
		Found:     true,
		Developer: toDeveloperItem(dev),
	}, nil
}
