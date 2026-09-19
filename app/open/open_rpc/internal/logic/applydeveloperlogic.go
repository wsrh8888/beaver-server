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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ApplyDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("apply_developer", ctx)}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(in *open_rpc.ApplyDeveloperReq) (*open_rpc.ApplyDeveloperRes, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "用户ID不能为空")
	}

	// 开发者登记不做审核：提交即生效，已有记录则覆盖登记信息。
	var existing open_models.OpenDeveloper
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&existing).Error
	if err == nil {
		existing.RealName = in.RealName
		existing.CompanyName = in.CompanyName
		existing.Phone = in.Phone
		existing.Email = in.Email
		if err := l.svcCtx.DB.Save(&existing).Error; err != nil {
			l.logger.Error(model.LogMsg{
				Text: "更新开发者登记失败",
				Data: map[string]interface{}{"err": err.Error()},
			})
			return nil, status.Error(codes.Internal, "提交失败")
		}
		return &open_rpc.ApplyDeveloperRes{Id: uint64(existing.Id)}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	dev := open_models.OpenDeveloper{
		UserID:      in.UserId,
		RealName:    in.RealName,
		CompanyName: in.CompanyName,
		Phone:       in.Phone,
		Email:       in.Email,
	}
	if err := l.svcCtx.DB.Create(&dev).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建开发者登记失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, status.Error(codes.Internal, "提交失败")
	}

	return &open_rpc.ApplyDeveloperRes{Id: uint64(dev.Id)}, nil
}
