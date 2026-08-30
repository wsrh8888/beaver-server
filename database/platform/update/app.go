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

package update

import (
	"beaver/app/platform/platform_models"
	"fmt"

	"gorm.io/gorm"
)

// DefaultUpdateAppID 海狸IM 默认升级应用 ID，与各端客户端保持一致
const DefaultUpdateAppID = "87c9dc499cc34f32896a4537e66cf65e"

// InitUpdateApp 幂等初始化升级默认应用
func InitUpdateApp(db *gorm.DB) error {
	fmt.Println("=== 开始初始化升级应用表 ===")

	defaultApp := platform_models.UpdateApp{
		Name:        "海狸IM",
		Description: "海狸IM",
		AppID:       DefaultUpdateAppID,
		Icon:        "",
		IsActive:    true,
	}

	var count int64
	if err := db.Model(&platform_models.UpdateApp{}).Where("app_id = ?", defaultApp.AppID).Count(&count).Error; err != nil {
		return fmt.Errorf("检查应用数据失败: %w", err)
	}

	if count == 0 {
		if err := db.Create(&defaultApp).Error; err != nil {
			return fmt.Errorf("创建应用数据失败: %w", err)
		}
		fmt.Printf("已创建应用: %s (AppID: %s)\n", defaultApp.Name, defaultApp.AppID)
	} else {
		fmt.Printf("应用已存在，跳过: %s (AppID: %s)\n", defaultApp.Name, defaultApp.AppID)
	}

	fmt.Println("=== 升级应用表初始化完成 ===")
	return nil
}
