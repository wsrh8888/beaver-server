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
	"beaver/app/platform/platform_models"

	"gorm.io/gorm"
)

func bucketNameMap(db *gorm.DB, bucketIDs []string) map[string]string {
	result := make(map[string]string, len(bucketIDs))
	if len(bucketIDs) == 0 {
		return result
	}

	uniqueIDs := make([]string, 0, len(bucketIDs))
	seen := make(map[string]struct{}, len(bucketIDs))
	for _, id := range bucketIDs {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return result
	}

	var buckets []platform_models.TrackBucket
	db.Where("bucket_id IN ?", uniqueIDs).Find(&buckets)
	for _, bucket := range buckets {
		result[bucket.BucketID] = bucket.Name
	}
	return result
}

func bucketNameByID(db *gorm.DB, bucketID string) string {
	if bucketID == "" {
		return ""
	}
	var bucket platform_models.TrackBucket
	if err := db.Where("bucket_id = ?", bucketID).Take(&bucket).Error; err != nil {
		return ""
	}
	return bucket.Name
}
