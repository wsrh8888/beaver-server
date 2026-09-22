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

// Organization 组织，表 org_organization。
// 一套安装只有一个组织，所以部门、员工不再重复存 org_id。
// 云上的飞书、钉钉会在每张表带 corp_id，那是多租户；私有化一套进程只服务这一家。
// 拥有者不放在这张表，见 Admin。
type Organization struct {
	models.Model
	OrgID string `gorm:"column:org_id;size:64;uniqueIndex;not null;comment:组织业务ID" json:"orgId"` // 组织 ID，唯一
	Name  string `gorm:"size:128;not null;comment:组织名" json:"name"`                              // 组织名，缺省「海狸」
}

func (Organization) TableName() string { return "org_organization" }
