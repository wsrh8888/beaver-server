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

package agent_models

import "beaver/common/models"

// 平台默认任务角色（对标 Agent 编排里的 workload routing）。
// 选型：先匹配能力，再按 tier / fallback，而不是只比价格。
const (
	ModelRoleRouter = "router" // 意图分类、是否调工具、链路选择 → 通常 fast
	ModelRoleChat   = "chat"   // 普通对话 → 通常 balanced
	ModelRoleTools  = "tools"  // 工具调用循环 → 要求 tools 能力
	ModelRoleVision = "vision" // 看图 → 要求 vision 能力
	ModelRoleReason = "reason" // 深推理 / 复杂规划 → 通常 strong
)

// AgentModelRole 平台「任务角色 → 官方模型」默认路由策略。
// 与模型目录分离：目录描述「能做什么」，本表描述「什么任务默认用谁」。
type AgentModelRole struct {
	models.Model
	Role            string `gorm:"size:32;uniqueIndex;not null;comment:任务角色router/chat/tools/vision/reason" json:"role"`
	ModelID         string `gorm:"size:64;index;not null;comment:主选官方模型ID" json:"modelId"`
	FallbackModelID string `gorm:"size:64;not null;default:'';comment:降级官方模型ID" json:"fallbackModelId"`
	Status          int8   `gorm:"type:tinyint;not null;default:1;index;comment:0停用1启用" json:"status"`
}

func (AgentModelRole) TableName() string {
	return "agent_model_roles"
}
