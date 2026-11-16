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

package feedback

import (
	"context"
	"errors"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SubmitFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("submit_feedback", ctx),
	}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(req *types.SubmitFeedbackReq) (*types.SubmitFeedbackRes, error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("反馈内容不能为空")
	}
	if req.Type < 1 || req.Type > 4 {
		return nil, errors.New("反馈类型不合法")
	}

	_, err := l.svcCtx.PlatformRpc.SubmitFeedback(l.ctx, &platform_rpc.SubmitFeedbackReq{
		UserId:    req.UserID,
		Content:   req.Content,
		Type:      int32(req.Type),
		FileNames: req.FileNames,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "提交反馈失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, errors.New("提交反馈失败")
	}

	l.logger.Info(model.LogMsg{Text: "提交反馈成功", Data: map[string]interface{}{"userId": req.UserID, "type": req.Type}})

	return &types.SubmitFeedbackRes{}, nil
}
