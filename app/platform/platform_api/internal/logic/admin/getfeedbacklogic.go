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

type GetFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：查询单条用户反馈
func NewGetFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackLogic {
	return &GetFeedbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetFeedback", ctx),
	}
}

// GetFeedback 只做协议转换，业务规则全部留在 platform_rpc。
func (l *GetFeedbackLogic) GetFeedback(req *types.GetFeedbackReq) (resp *types.GetFeedbackRes, err error) {
	rpcReq := &platformcli.GetFeedbackReq{Id: req.ID}
	rpcRes, err := l.svcCtx.PlatformRpc.GetFeedback(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.GetFeedbackRes{Feedback: &types.FeedbackItem{ID: rpcRes.Feedback.Id, UserID: rpcRes.Feedback.UserId, Content: rpcRes.Feedback.Content, Type: rpcRes.Feedback.Type, Status: rpcRes.Feedback.Status, FileNames: rpcRes.Feedback.FileNames, HandlerID: rpcRes.Feedback.HandlerId, HandleTime: rpcRes.Feedback.HandleTime, HandleResult: rpcRes.Feedback.HandleResult, CreatedAt: rpcRes.Feedback.CreatedAt, UpdatedAt: rpcRes.Feedback.UpdatedAt}}, nil
}
