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

package workbench

import (
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toProtoEntryConfig(cfg *types.WorkbenchEntryConfig) *platform_rpc.WorkbenchEntryConfig {
	if cfg == nil {
		return nil
	}
	return &platform_rpc.WorkbenchEntryConfig{
		Type:   int32(cfg.Type),
		Pc:     cfg.PC,
		Mobile: cfg.Mobile,
	}
}

func fromProtoEntryConfig(cfg *platform_rpc.WorkbenchEntryConfig) types.WorkbenchEntryConfig {
	if cfg == nil {
		return types.WorkbenchEntryConfig{}
	}
	return types.WorkbenchEntryConfig{
		Type:   int(cfg.Type),
		PC:     cfg.Pc,
		Mobile: cfg.Mobile,
	}
}

func toAdminAppItem(item *platform_rpc.WorkbenchAppItem) types.GetWorkbenchAppListItem {
	return types.GetWorkbenchAppListItem{
		WorkbenchAppID: item.WorkbenchAppId,
		Name:           item.Name,
		Description:    item.Description,
		Icon:           item.Icon,
		AppType:        int(item.AppType),
		ClientScope:    int(item.ClientScope),
		EntryConfig:    fromProtoEntryConfig(item.EntryConfig),
		OpenMode:       int(item.OpenMode),
		Category:       int(item.Category),
		Sort:           int(item.Sort),
		Status:         int(item.Status),
		Remark:         item.Remark,
		CreatedBy:      item.CreatedBy,
		LastModifiedBy: item.LastModifiedBy,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}
