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

package admin

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	platformcli "beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
)

type ListContentReportsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询内容举报
func NewListContentReportsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContentReportsLogic {
	return &ListContentReportsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListContentReports", ctx),
	}
}

// ListContentReports 只做协议转换，业务规则全部留在 platform_rpc。
func (l *ListContentReportsLogic) ListContentReports(req *types.ListContentReportsReq) (resp *types.ListContentReportsRes, err error) {
	rpcReq := &platformcli.ListContentReportsReq{Page: req.Page, PageSize: req.PageSize, Status: req.Status, TargetType: req.TargetType, TargetId: req.TargetID, ReporterUserId: req.ReporterUserID}
	rpcRes, err := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.ContentReportItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.ContentReportItem{ID: it1.Id, ReporterUserID: it1.ReporterUserId, TargetType: it1.TargetType, TargetID: it1.TargetId, ReasonType: it1.ReasonType, Content: it1.Content, FileNames: it1.FileNames, Status: it1.Status, CaseID: it1.CaseId, HandlerID: it1.HandlerId, HandleRemark: it1.HandleRemark, CreatedAt: it1.CreatedAt})
	}
	return &types.ListContentReportsRes{Total: rpcRes.Total, List: tmp2}, nil
}
