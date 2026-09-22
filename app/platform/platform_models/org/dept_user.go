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

// DeptUser 人与部门的关系，表 org_dept_user。
// 一个人可以属于多个部门，其中一行 IsMain 为主部门。同一人同一部门只有一行。
// IsLeader 表示是这个部门的负责人，一个部门可以有多个负责人。企业微信的 is_leader_in_dept 就是这个。
type DeptUser struct {
	models.Model
	UserID   string `gorm:"column:user_id;size:64;not null;uniqueIndex:uk_user_dept;index;comment:用户ID" json:"userId"` // 对应 Staff.UserID
	DeptID   string `gorm:"column:dept_id;size:64;not null;uniqueIndex:uk_user_dept;index;comment:部门ID" json:"deptId"` // 对应 Department.DeptID
	IsMain   bool   `gorm:"not null;default:false;comment:是否主部门" json:"isMain"`                                        // 主部门。一个人只有一个主部门
	IsLeader bool   `gorm:"not null;default:false;comment:是否该部门负责人" json:"isLeader"`                                   // 是否这个部门的负责人
	Sort     int    `gorm:"not null;default:0;comment:部门内排序" json:"sort"`                                              // 在这个部门里的排序，越小越靠前
}

func (DeptUser) TableName() string { return "org_dept_user" }
