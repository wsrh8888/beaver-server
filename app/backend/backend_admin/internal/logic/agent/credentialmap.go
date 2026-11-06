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

// toAgentCredentialInfo 接入通道对外项映射。platform 侧已只给 key_hint/key_configured，此处原样透出。
func toAgentCredentialInfo(item *platform_rpc.AgentCredentialItem) types.AgentCredentialInfo {
	if item == nil {
		return types.AgentCredentialInfo{}
	}
	return types.AgentCredentialInfo{
		CredentialID:  item.CredentialId,
		Name:          item.Name,
		ModelService:  item.ModelService,
		ChannelType:   item.ChannelType,
		BaseURL:       item.BaseUrl,
		KeyHint:       item.KeyHint,
		KeyConfigured: item.KeyConfigured,
		Sort:          int(item.Sort),
		Status:        int(item.Status),
	}
}
