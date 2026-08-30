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

package platform_models

import "beaver/common/models"

// TrackBucket 埋点与日志的 Bucket 注册表
type TrackBucket struct {
	models.Model
	Name        string `json:"name" gorm:"size:64"`
	Description string `json:"description" gorm:"type:text"`
	BucketID    string `json:"bucketId" gorm:"column:bucket_id;uniqueIndex;size:64"`
	Kind        string `json:"kind" gorm:"index;size:16"` // track=埋点 log=日志
	CreateUser  string `json:"createUser" gorm:"size:64"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
}
