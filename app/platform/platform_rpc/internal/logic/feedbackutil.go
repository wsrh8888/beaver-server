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

func toFeedbackItem(f platform_models.FeedbackModel) *platform_rpc.FeedbackItem {
	handleTime := ""
	if f.HandleTime != nil {
		handleTime = f.HandleTime.Format(time.RFC3339)
	}
	return &platform_rpc.FeedbackItem{
		Id:           uint64(f.Id),
		UserId:       f.UserID,
		Content:      f.Content,
		Type:         int32(f.Type),
		Status:       int32(f.Status),
		FileNames:    []string(f.FileNames),
		HandlerId:    f.HandlerID,
		HandleTime:   handleTime,
		HandleResult: f.HandleResult,
		CreatedAt:    time.Time(f.CreatedAt).Format(time.RFC3339),
		UpdatedAt:    time.Time(f.UpdatedAt).Format(time.RFC3339),
	}
}
