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

package track_public

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LogEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogEventsLogic {
	return &LogEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogEventsLogic) LogEvents(req *types.LogEventsReq) (*types.LogEventsRes, error) {
	logs := make([]platform_models.TrackLogger, 0, len(req.Logs))
	for _, logData := range req.Logs {
		logs = append(logs, platform_models.TrackLogger{
			Level:     logData.Level,
			Data:      datatypes.JSON(logData.Data),
			BucketID:  logData.BucketID,
			Timestamp: logData.Timestamp,
		})
	}

	go func(logsToSave []platform_models.TrackLogger) {
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(logsToSave, 100).Error
		})
		if err != nil {
			l.Logger.Errorf("save logs failed: %v", err)
		}
	}(logs)

	return &types.LogEventsRes{}, nil
}
