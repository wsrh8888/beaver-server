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

// OfficialModelIDAuto 虚拟官方模型：客户端默认选中；agent 收到后自行路由真实模型。
const OfficialModelIDAuto = "auto"

// AgentOfficialModel 平台官方模型（全员可读，无 UserID）。
// 与自定义同构：Tier + Capabilities + Endpoint + ApiKey；独有：Provider / ModelName / Sort。
// ApiKey 为平台侧密钥（后续加密）；本地可种子占位后手动改库。
type AgentOfficialModel struct {
	models.Model
	ModelID      string                 `gorm:"size:64;uniqueIndex;not null;comment:官方模型业务ID" json:"modelId"`
	Name         string                 `gorm:"size:128;not null;comment:展示名称" json:"name"`
	Provider     string                 `gorm:"size:64;not null;default:'';comment:供应商标识" json:"provider"`
	ModelName    string                 `gorm:"size:128;not null;comment:实际调用模型名" json:"modelName"`
	Endpoint     string                 `gorm:"size:500;not null;default:'';comment:接口地址(空则平台默认网关)" json:"endpoint"`
	ApiKey       string                 `gorm:"size:512;not null;default:'';comment:平台API Key(后续加密存储)" json:"apiKey"`
	Tier         string                 `gorm:"size:16;not null;default:balanced;index;comment:档位fast/balanced/strong" json:"tier"`
	Capabilities AgentModelCapabilities `gorm:"type:longtext;comment:能力画像JSON(vision/tools/reasoning等)" json:"capabilities"`
	Sort         int32                  `gorm:"not null;default:0;index;comment:排序越小越靠前" json:"sort"`
	Status       int8                   `gorm:"type:tinyint;not null;default:1;index;comment:0下架1上架" json:"status"`
}

func (AgentOfficialModel) TableName() string {
	return "agent_official_models"
}
