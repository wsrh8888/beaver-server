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
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

// validateKnowledgeBindings 保存 AgentApp 时校验其 knowledge_ids 绑定：
// 逐个回查 platform_rpc，集合不存在即拒绝——防止把应用绑到不存在/已删的知识库，
// 让运行时的知识检索在配置期就对齐，而非等到发消息才发现查不到。
// knowledge_ids 为空或解析失败视为「未绑定」，放行（向后兼容旧配置）。
func validateKnowledgeBindings(ctx context.Context, svcCtx *svc.ServiceContext, knowledgeIDsJSON string) error {
	raw := strings.TrimSpace(knowledgeIDsJSON)
	if raw == "" {
		return nil
	}
	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil // 坏 JSON 不当作绑定意图，交由既有配置流程处理
	}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		res, err := svcCtx.PlatformRpc.GetCollection(ctx, &platform_rpc.GetCollectionReq{CollectionId: id})
		if err != nil {
			return fmt.Errorf("校验知识库绑定失败 %s: %w", id, err)
		}
		if res.GetCollection() == nil || res.GetCollection().GetCollectionId() == "" {
			return fmt.Errorf("知识库集合不存在: %s", id)
		}
	}
	return nil
}
