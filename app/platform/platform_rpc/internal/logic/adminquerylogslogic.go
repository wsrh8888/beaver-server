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
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminQueryLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminQueryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminQueryLogsLogic {
	return &AdminQueryLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminQueryLogsLogic) AdminQueryLogs(in *platform_rpc.AdminQueryLogsReq) (*platform_rpc.AdminQueryLogsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&platform_models.TrackLogger{})
	db = db.Where("bucket_id = ?", in.BucketId)
	if in.Level != "" {
		db = db.Where("level = ?", in.Level)
	}
	if in.Keyword != "" {
		db = db.Where("data LIKE ?", "%"+in.Keyword+"%")
	}
	if in.StartTime > 0 {
		db = db.Where("timestamp >= ?", in.StartTime)
	}
	if in.EndTime > 0 {
		db = db.Where("timestamp <= ?", in.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("查询日志总数失败: %v", err)
		return nil, err
	}

	var logs []platform_models.TrackLogger
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("timestamp DESC").Find(&logs).Error; err != nil {
		l.Errorf("查询日志列表失败: %v", err)
		return nil, err
	}

	bucketName := bucketNameByID(l.svcCtx.DB, in.BucketId)

	items := make([]*platform_rpc.AdminLogItem, 0, len(logs))
	for _, logItem := range logs {
		dataStr := ""
		if logItem.Data != nil {
			dataBytes, _ := logItem.Data.MarshalJSON()
			dataStr = string(dataBytes)
		}
		items = append(items, &platform_rpc.AdminLogItem{
			Id:         uint64(logItem.Id),
			Level:      logItem.Level,
			Data:       dataStr,
			BucketId:   logItem.BucketID,
			BucketName: bucketName,
			Timestamp:  logItem.Timestamp,
			CreatedAt:  time.Time(logItem.CreatedAt).Format(time.RFC3339),
		})
	}

	return &platform_rpc.AdminQueryLogsRes{Total: total, Logs: items}, nil
}
