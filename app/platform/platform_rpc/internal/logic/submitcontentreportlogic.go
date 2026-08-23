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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmitContentReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitContentReportLogic {
	return &SubmitContentReportLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SubmitContentReportLogic) SubmitContentReport(in *platform_rpc.SubmitContentReportReq) (*platform_rpc.SubmitContentReportRes, error) {
	if in.ReporterUserId == "" || in.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "举报人或举报对象不能为空")
	}
	report := platform_models.ContentReportModel{
		ReporterUserID: in.ReporterUserId,
		TargetType:     int(in.TargetType),
		TargetID:       in.TargetId,
		ReasonType:     int(in.ReasonType),
		Content:        in.Content,
		FileNames:      platform_models.FileNames(in.FileNames),
		Status:         platform_models.ReportStatusPending,
	}
	if err := l.svcCtx.DB.Create(&report).Error; err != nil {
		l.Errorf("创建举报失败: %v", err)
		return nil, status.Error(codes.Internal, "提交失败")
	}
	return &platform_rpc.SubmitContentReportRes{Id: uint64(report.Id)}, nil
}
