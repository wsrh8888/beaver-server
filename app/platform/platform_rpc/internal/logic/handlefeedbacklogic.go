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
	"strconv"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type HandleFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewHandleFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFeedbackLogic {
	return &HandleFeedbackLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("handle_feedback", ctx)}
}

func (l *HandleFeedbackLogic) HandleFeedback(in *platform_rpc.HandleFeedbackReq) (*platform_rpc.HandleFeedbackRes, error) {
	if in.Status < 1 || in.Status > 4 {
		return nil, status.Error(codes.InvalidArgument, "无效的状态值")
	}

	var feedback platform_models.FeedbackModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "反馈记录不存在")
		}
		l.logger.Error(model.LogMsg{Text: "查询反馈失败", Data: map[string]any{"id": in.Id, "err": err.Error()}})
		return nil, err
	}

	handlerID, _ := strconv.ParseInt(in.HandlerId, 10, 64)
	now := time.Now()
	if err := l.svcCtx.DB.Model(&feedback).Updates(map[string]interface{}{
		"status":        platform_models.FeedbackStatus(in.Status),
		"handle_result": in.HandleResult,
		"handler_id":    handlerID,
		"handle_time":   &now,
		"updated_at":    now,
	}).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "处理反馈失败", Data: map[string]any{"id": in.Id, "err": err.Error()}})
		return nil, status.Error(codes.Internal, "处理反馈失败")
	}

	l.logger.Info(model.LogMsg{Text: "处理反馈成功", Data: map[string]interface{}{"id": in.Id, "status": in.Status}})

	return &platform_rpc.HandleFeedbackRes{}, nil
}
