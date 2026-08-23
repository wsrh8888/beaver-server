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

package backend_models

import (
	"beaver/common/models"
)

/**
 * @description: 存储系统中的角色信息
 * 注意：微服务架构中不使用数据库外键约束，关联关系通过关联表维护
 */
type AdminSystemAuthority struct {
	models.Model
	Name        string `json:"authorityName" gorm:"size:32;unique;index;comment:角色名"` // 角色名（唯一）
	Description string `json:"description" gorm:"size:256;comment:角色描述"`              // 角色描述
	Status      int8   `json:"status" gorm:"default:1;index;comment:状态"`              // 1:启用 2:禁用
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`                      // 排序
}
