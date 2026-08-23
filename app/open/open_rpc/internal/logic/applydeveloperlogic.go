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

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ApplyDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(in *open_rpc.ApplyDeveloperReq) (*open_rpc.ApplyDeveloperRes, error) {
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "用户ID不能为空")
	}

	var existing open_models.OpenDeveloper
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&existing).Error
	if err == nil {
		switch existing.Status {
		case 1:
			return nil, status.Error(codes.AlreadyExists, "您已经是认证开发者")
		case 0:
			return nil, status.Error(codes.AlreadyExists, "申请正在审核中")
		case 2:
			existing.RealName = in.RealName
			existing.CompanyName = in.CompanyName
			existing.Phone = in.Phone
			existing.Email = in.Email
			existing.Description = in.Description
			existing.Status = 0
			existing.AuditBy = ""
			existing.AuditTime = 0
			existing.AuditRemark = ""
			if err := l.svcCtx.DB.Save(&existing).Error; err != nil {
				l.Errorf("重新申请开发者失败: %v", err)
				return nil, status.Error(codes.Internal, "申请失败")
			}
			return &open_rpc.ApplyDeveloperRes{Id: uint64(existing.Id)}, nil
		default:
			return nil, status.Error(codes.FailedPrecondition, "申请状态异常")
		}
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
		Description: in.Description,
		Status:      0,
	}
	if err := l.svcCtx.DB.Create(&dev).Error; err != nil {
		l.Errorf("创建开发者申请失败: %v", err)
		return nil, status.Error(codes.Internal, "申请失败")
	}

	return &open_rpc.ApplyDeveloperRes{Id: uint64(dev.Id)}, nil
}
