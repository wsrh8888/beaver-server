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

package org_models

import "beaver/common/models"

const (
	AdminOwner int8 = 1 // 创建者。全组织只有一个，不能移出
	AdminRole  int8 = 2 // 组织管理员。可以管理部门和成员，不是某个部门的负责人
)

// Admin 组织管理员，表 org_admin。
// 普通员工不在这张表里。部门负责人看 DeptUser.IsLeader，不要和这里的管理员混在一起。
type Admin struct {
	models.Model
	UserID string `gorm:"column:user_id;size:64;uniqueIndex;not null;comment:用户ID" json:"userId"` // 对应 Staff.UserID，一个人只有一种管理身份
	Role   int8   `gorm:"type:tinyint;not null;comment:1创建者 2组织管理员" json:"role"`                  // 1 创建者，2 组织管理员
}

func (Admin) TableName() string { return "org_admin" }
