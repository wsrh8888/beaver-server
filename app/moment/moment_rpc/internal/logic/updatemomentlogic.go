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

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type UpdateMomentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMomentLogic {
	return &UpdateMomentLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("update_moment", ctx)}
}

func (l *UpdateMomentLogic) UpdateMoment(in *moment_rpc.UpdateMomentReq) (*moment_rpc.UpdateMomentRes, error) {
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.Where("moment_id = ?", in.MomentId).First(&moment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("动态不存在")
		}
		l.logger.Error(model.LogMsg{
			Text: "查询动态失败",
			Data: map[string]any{"momentId": in.MomentId, "err": err.Error()},
		})
		return nil, err
	}

	if in.IsDeleted != nil && *in.IsDeleted {
		if err := l.svcCtx.DB.Model(&moment).Update("is_deleted", true).Error; err != nil {
			l.logger.Error(model.LogMsg{
				Text: "更新动态失败",
				Data: map[string]any{"momentId": in.MomentId, "err": err.Error()},
			})
			return nil, err
		}
		_ = l.svcCtx.DB.Where("moment_id = ?", in.MomentId).Delete(&moment_models.MomentCommentModel{}).Error
		_ = l.svcCtx.DB.Where("moment_id = ?", in.MomentId).Delete(&moment_models.MomentLikeModel{}).Error
	}

	l.logger.Info(model.LogMsg{
		Text: "更新动态成功",
		Data: map[string]any{"momentId": in.MomentId, "isDeleted": in.IsDeleted},
	})
	return &moment_rpc.UpdateMomentRes{}, nil
}
