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
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CreateAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("create_app", ctx)}
}

func (l *CreateAppLogic) CreateApp(in *platform_rpc.CreateAppReq) (*platform_rpc.CreateAppRes, error) {
	var existing platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("name = ?", in.Name).First(&existing).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "应用名称已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.logger.Error(model.LogMsg{Text: "查询应用名称唯一性失败", Data: map[string]any{"name": in.Name, "err": err.Error()}})
		return nil, err
	}

	app := platform_models.UpdateApp{
		Name:        in.Name,
		AppID:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		Description: in.Description,
		IsActive:    true,
	}
	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "创建应用失败", Data: map[string]any{"name": in.Name, "err": err.Error()}})
		return nil, status.Error(codes.Internal, "创建应用失败")
	}

	l.logger.Info(model.LogMsg{Text: "创建应用成功", Data: map[string]interface{}{"appId": app.AppID, "name": in.Name}})

	return &platform_rpc.CreateAppRes{Id: uint64(app.Id), AppId: app.AppID}, nil
}
