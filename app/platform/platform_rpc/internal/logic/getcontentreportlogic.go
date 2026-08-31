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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetContentReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentReportLogic {
	return &GetContentReportLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("get_content_report", ctx)}
}

func (l *GetContentReportLogic) GetContentReport(in *platform_rpc.GetContentReportReq) (*platform_rpc.GetContentReportRes, error) {
	var report platform_models.ContentReportModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "举报不存在")
		}
		l.logger.Error(model.LogMsg{Text: "查询举报失败", Data: map[string]any{"id": in.Id, "err": err.Error()}})
		return nil, err
	}
	return &platform_rpc.GetContentReportRes{Report: toContentReportItem(report)}, nil
}
