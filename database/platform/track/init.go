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

package track

import (
	"beaver/app/platform/platform_models"
	"fmt"

	"gorm.io/gorm"
)

// InitBuckets 幂等初始化默认 track_buckets（日志 + 埋点）
func InitBuckets(db *gorm.DB) error {
	buckets := []platform_models.TrackBucket{
		{
			Name:        "海狸客户端-日志",
			Description: "海狸 IM 客户端运行日志（log）",
			BucketID:    "b2c3d4e5-f6a7-4789-b012-456789abcdef",
			Kind:        "log",
			CreateUser:  "system",
			IsActive:    true,
		},
		{
			Name:        "海狸客户端-埋点",
			Description: "海狸 IM 客户端行为埋点（report_events）",
			BucketID:    "a1b2c3d4-e5f6-4789-a012-3456789abcde",
			Kind:        "track",
			CreateUser:  "system",
			IsActive:    true,
		},
	}

	fmt.Println("=== 开始初始化 track_buckets 默认数据 ===")
	for _, bucket := range buckets {
		if err := upsertBucket(db, bucket); err != nil {
			return err
		}
	}
	fmt.Println("=== track_buckets 默认数据初始化完成 ===")
	return nil
}

func upsertBucket(db *gorm.DB, bucket platform_models.TrackBucket) error {
	var count int64
	if err := db.Model(&platform_models.TrackBucket{}).Where("bucket_id = ?", bucket.BucketID).Count(&count).Error; err != nil {
		return fmt.Errorf("检查 Track Bucket 失败: %w", err)
	}
	if count == 0 {
		if err := db.Create(&bucket).Error; err != nil {
			return fmt.Errorf("创建 Track Bucket 失败: %w", err)
		}
		fmt.Printf("已创建 Track Bucket: %s (BucketID: %s, Kind: %s)\n", bucket.Name, bucket.BucketID, bucket.Kind)
		return nil
	}

	if err := db.Model(&platform_models.TrackBucket{}).Where("bucket_id = ?", bucket.BucketID).Updates(map[string]any{
		"name":        bucket.Name,
		"description": bucket.Description,
		"kind":        bucket.Kind,
		"is_active":   bucket.IsActive,
	}).Error; err != nil {
		return fmt.Errorf("更新 Track Bucket 失败: %w", err)
	}
	fmt.Printf("Track Bucket 已存在，已同步: %s (BucketID: %s, Kind: %s)\n", bucket.Name, bucket.BucketID, bucket.Kind)
	return nil
}
