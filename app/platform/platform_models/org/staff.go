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
	StaffActive   int8 = 1 // 在职。对应飞书已激活、企业微信已激活
	StaffFrozen   int8 = 2 // 停用。账号还在，不能正常使用
	StaffInactive int8 = 3 // 未激活。已拉进组织，人还没进来
	StaffResigned int8 = 4 // 离职。记录保留，通讯录默认不展示
)

// Staff 员工档案，表 org_staff。一个人在这个组织里一行。
// 这是通讯录里的人，不是部门关系。昵称、头像仍在用户表。
// 直属上级是汇报线，和「某部门的负责人」不是一回事，负责人在 DeptUser.IsLeader。
type Staff struct {
	models.Model
	UserID       string `gorm:"column:user_id;size:64;uniqueIndex;not null;comment:用户ID" json:"userId"`                       // 对应 UserModel.UserID
	EmployeeNo   string `gorm:"column:employee_no;size:64;not null;default:'';comment:工号" json:"employeeNo"`                  // 工号，可空
	JobTitle     string `gorm:"column:job_title;size:128;not null;default:'';comment:职位" json:"jobTitle"`                     // 职位，通讯录卡片上的 title
	LeaderUserID string `gorm:"column:leader_user_id;size:64;not null;default:'';index;comment:直属上级用户ID" json:"leaderUserId"` // 直属上级，空字符串表示没有
	Status       int8   `gorm:"type:tinyint;not null;default:1;index;comment:1在职 2停用 3未激活 4离职" json:"status"`                 // 1 在职，2 停用，3 未激活，4 离职
}

func (Staff) TableName() string { return "org_staff" }
