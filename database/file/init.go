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

package fileseed

import (
	"beaver/app/file/file_models"
	"fmt"

	"gorm.io/gorm"
)

// InitDefaultFiles 初始化默认文件数据
func InitDefaultFiles(db *gorm.DB) error {
	defaultFiles := []file_models.FileModel{
		{
			OriginalName: "defaultUserFileName",
			Size:         60317,
			Path:         "image/user.png",
			Md5:          "a9de5548bef8c10b92428fff61275c72",
			Type:         "image",
			FileKey:      "a9de5548bef8c10b92428fff61275c72.png",
			Source:       file_models.QiniuSource,
		},
		{
			OriginalName: "defaultGroupFileName",
			Size:         90310,
			Path:         "image/group.png",
			Md5:          "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Type:         "image",
			FileKey:      "a8ba5d19ea54a91aec17dec0ad5000e6.png",
			Source:       file_models.QiniuSource,
		},
	}

	for _, file := range defaultFiles {
		var count int64
		if err := db.Model(&file_models.FileModel{}).Where("file_key = ?", file.FileKey).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&file).Error; err != nil {
				return fmt.Errorf("创建默认文件失败: %w", err)
			}
			fmt.Printf("创建默认文件成功: %s\n", file.FileKey)
		}
	}

	return nil
}
