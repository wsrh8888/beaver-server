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

package open_models

import "beaver/common/models"

// OpenDeveloper 开发者登记（一个用户一份开发者身份）
//
// 不做审核：登记即可创建应用 —— "谁有资格创建应用"由这条记录本身表达。
// 因此不存任何流程态字段（待审核 / 审核人 / 审核时间 / 申请说明）。
type OpenDeveloper struct {
	models.Model
	UserID      string `gorm:"size:64;uniqueIndex;not null;comment:用户ID"`
	RealName    string `gorm:"size:32;comment:真实姓名"`
	CompanyName string `gorm:"size:64;comment:公司名称"`
	Phone       string `gorm:"size:11;comment:联系电话"`
	Email       string `gorm:"size:128;comment:邮箱"`
}
