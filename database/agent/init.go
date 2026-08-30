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

package agentseed

import (
	"beaver/app/agent/agent_models"
	"fmt"

	"gorm.io/gorm"
)

const defaultDeepSeekModelID = "deepseek-chat"

// InitOfficialModels 初始化官方模型目录（幂等：已存在则跳过，不覆盖手动改的 key）。
func InitOfficialModels(db *gorm.DB) error {
	row := agent_models.AgentOfficialModel{
		ModelID:   defaultDeepSeekModelID,
		Name:      "DeepSeek Chat",
		Provider:  "deepseek",
		ModelName: "deepseek-chat",
		Endpoint:  "https://api.deepseek.com",
		ApiKey:    "", // 本地种子占位，启动后手动 UPDATE api_key
		Tier:      agent_models.ModelTierBalanced,
		Capabilities: agent_models.AgentModelCapabilities{
			Vision:           false,
			Tools:            true,
			Reasoning:        true,
			StructuredOutput: true,
			LongContext:      true,
		},
		Sort:   0,
		Status: 1,
	}

	var count int64
	if err := db.Model(&agent_models.AgentOfficialModel{}).
		Where("model_id = ?", row.ModelID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if err := db.Create(&row).Error; err != nil {
		return fmt.Errorf("创建默认官方模型失败: %w", err)
	}
	fmt.Printf("创建默认官方模型成功: %s（请手动填写 api_key）\n", row.ModelID)
	return nil
}
