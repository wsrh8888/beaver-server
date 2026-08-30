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

import "beaver/common/models"

/**
 * @description: 表示用户-角色关联关系数据模型
 * 注意：微服务架构中不使用数据库外键约束，数据一致性由应用层保证
 */
type AdminSystemAuthorityUser struct {
	models.Model
	UserID      string `json:"userId" gorm:"size:64;comment:用户ID;index:idx_user_authority,unique"`
	AuthorityID uint   `json:"authorityId" gorm:"comment:角色ID;index:idx_user_authority,unique"`
}
