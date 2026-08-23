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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateArchitectureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArchitectureLogic {
	return &UpdateArchitectureLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateArchitectureLogic) UpdateArchitecture(in *platform_rpc.UpdateArchitectureReq) (*platform_rpc.UpdateArchitectureRes, error) {
	var arch platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.First(&arch, in.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "架构不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{"is_active": in.IsActive}
	if in.Description != "" {
		updates["description"] = in.Description
	}
	if err := l.svcCtx.DB.Model(&arch).Updates(updates).Error; err != nil {
		l.Errorf("更新架构失败: %v", err)
		return nil, status.Error(codes.Internal, "更新架构失败")
	}

	return &platform_rpc.UpdateArchitectureRes{}, nil
}
