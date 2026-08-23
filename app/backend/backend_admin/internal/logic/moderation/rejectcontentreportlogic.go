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

package moderation

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectContentReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectContentReportLogic {
	return &RejectContentReportLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RejectContentReportLogic) RejectContentReport(req *types.RejectContentReportReq) (resp *types.RejectContentReportRes, err error) {
	if req.ReportID == 0 {
		return nil, errors.New("举报ID不能为空")
	}

	reportRes, err := l.svcCtx.PlatformRpc.GetContentReport(l.ctx, &platform_rpc.GetContentReportReq{Id: req.ReportID})
	if err != nil {
		l.Errorf("获取举报详情失败: %v", err)
		return nil, err
	}
	if reportRes.Report == nil {
		return nil, errors.New("举报不存在")
	}
	report := reportRes.Report
	if report.Status != platform_models.ReportStatusPending {
		return nil, errors.New("仅待处理举报可驳回")
	}

	remark := req.HandleRemark
	if remark == "" {
		remark = "举报驳回"
	}

	_, err = l.svcCtx.PlatformRpc.UpdateContentReports(l.ctx, &platform_rpc.UpdateContentReportsReq{
		Ids:          []uint64{report.Id},
		Action:       2,
		HandlerId:    req.UserID,
		HandleRemark: remark,
	})
	if err != nil {
		l.Errorf("驳回举报失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "reject_report", "report", fmt.Sprintf("%d", report.Id), 0,
		fmt.Sprintf("驳回举报 targetType=%d targetId=%s remark=%s", report.TargetType, report.TargetId, remark), "success", "")

	return &types.RejectContentReportRes{}, nil
}
