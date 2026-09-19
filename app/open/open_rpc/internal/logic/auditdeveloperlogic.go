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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuditDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("audit_developer", ctx)}
}

// AuditDeveloper 开发者登记已改为「登记即生效」，不再有审核流程。
// 该接口保留仅为兼容既有调用方：只校验记录存在，不做任何状态变更。
func (l *AuditDeveloperLogic) AuditDeveloper(in *open_rpc.AuditDeveloperReq) (*open_rpc.AuditDeveloperRes, error) {
	var dev open_models.OpenDeveloper
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&dev).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "开发者记录不存在")
		}
		return nil, err
	}

	return &open_rpc.AuditDeveloperRes{}, nil
}
