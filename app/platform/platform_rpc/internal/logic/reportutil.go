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
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toContentReportItem(r platform_models.ContentReportModel) *platform_rpc.ContentReportItem {
	return &platform_rpc.ContentReportItem{
		Id:             uint64(r.Id),
		ReporterUserId: r.ReporterUserID,
		TargetType:     int32(r.TargetType),
		TargetId:       r.TargetID,
		ReasonType:     int32(r.ReasonType),
		Content:        r.Content,
		FileNames:      []string(r.FileNames),
		Status:         int32(r.Status),
		CaseId:         r.CaseID,
		HandlerId:      r.HandlerID,
		HandleRemark:   r.HandleRemark,
		CreatedAt:      time.Time(r.CreatedAt).Format(time.RFC3339),
	}
}
