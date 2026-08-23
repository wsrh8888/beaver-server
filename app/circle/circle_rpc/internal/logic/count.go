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

package logic

import (
	"beaver/app/circle/circle_models"

	"gorm.io/gorm"
)

type circleCountRow struct {
	CircleID string
	Count    int64
}

func countMembersByCircleIDs(db *gorm.DB, circleIDs []string) map[string]int64 {
	result := make(map[string]int64, len(circleIDs))
	if len(circleIDs) == 0 {
		return result
	}
	var rows []circleCountRow
	db.Model(&circle_models.CircleMemberModel{}).
		Select("circle_id, count(*) as count").
		Where("circle_id IN ?", circleIDs).
		Group("circle_id").
		Scan(&rows)
	for _, row := range rows {
		result[row.CircleID] = row.Count
	}
	return result
}

func countPostsByCircleIDs(db *gorm.DB, circleIDs []string) map[string]int64 {
	result := make(map[string]int64, len(circleIDs))
	if len(circleIDs) == 0 {
		return result
	}
	var rows []circleCountRow
	db.Model(&circle_models.CirclePostModel{}).
		Select("circle_id, count(*) as count").
		Where("circle_id IN ? AND is_deleted = false", circleIDs).
		Group("circle_id").
		Scan(&rows)
	for _, row := range rows {
		result[row.CircleID] = row.Count
	}
	return result
}
