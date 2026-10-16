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

package agent

import (
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toAgentModelInfo(item *platform_rpc.AgentModelItem) types.AgentModelInfo {
	if item == nil {
		return types.AgentModelInfo{}
	}
	return types.AgentModelInfo{
		ModelID:         item.ModelId,
		Source:          item.Source,
		UserID:          item.UserId,
		Name:            item.Name,
		Provider:        item.Provider,
		ModelName:       item.ModelName,
		Endpoint:        item.Endpoint,
		APIKey:          maskAPIKey(item.ApiKey),
		ContextWindow:   item.ContextWindow,
		MaxOutputTokens: int(item.MaxOutputTokens),
		SupportsTools:   item.SupportsTools,
		SupportsVision:  item.SupportsVision,
		Sort:            int(item.Sort),
		Status:          int(item.Status),
	}
}

func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
