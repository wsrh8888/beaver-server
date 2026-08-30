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

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetEventVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEventVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventVersionsLogic {
	return &GetEventVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEventVersionsLogic) GetEventVersions(in *notification_rpc.GetEventVersionsReq) (*notification_rpc.GetEventVersionsRes, error) {
	resp := &notification_rpc.GetEventVersionsRes{
		EventVersions: []*notification_rpc.EventVersion{},
		MaxVersion:    0,
	}

	var rows []notification_models.NotificationEvent
	query := l.svcCtx.DB.WithContext(l.ctx).
		Where("version > ?", in.SinceVersion).
		Order("version ASC")

	if in.Limit > 0 {
		query = query.Limit(int(in.Limit))
	}

	if err := query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		resp.EventVersions = append(resp.EventVersions, &notification_rpc.EventVersion{
			EventId: row.EventID,
			Version: row.Version,
		})
		if row.Version > resp.MaxVersion {
			resp.MaxVersion = row.Version
		}
	}

	return resp, nil
}
