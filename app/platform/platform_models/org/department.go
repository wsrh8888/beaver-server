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

// Department 部门，表 org_department。
// ParentID 为空字符串表示根。根只有一个，不可删除，也不能改父级。
// Path 是从根到自身的路径，查子树用前缀，避免每次递归。飞书、钉钉内部部门树都是这么查的。
type Department struct {
	models.Model
	DeptID   string `gorm:"column:dept_id;size:64;uniqueIndex;not null;comment:部门业务ID" json:"deptId"`               // 部门 ID，唯一
	ParentID string `gorm:"column:parent_id;size:64;not null;default:'';index;comment:父部门ID空字符串为根" json:"parentId"` // 父部门 ID，空字符串表示根
	Name     string `gorm:"size:128;not null;comment:部门名" json:"name"`                                              // 部门名
	Sort     int    `gorm:"not null;default:0;comment:同级排序" json:"sort"`                                            // 同级排序，越小越靠前
	Path     string `gorm:"size:1024;not null;default:'';index;comment:从根到自身的路径" json:"path"`                       // 如 /根部门ID/当前部门ID
}

func (Department) TableName() string { return "org_department" }
