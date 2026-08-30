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

package circle_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PostFileInfo 帖子附件信息
type PostFileInfo struct {
	FileKey string `json:"fileKey"` // 文件key
	Type    uint32 `json:"type"`    // 文件类型：2=图片 3=视频 4=文件
}

type PostFiles []PostFileInfo

func (f *PostFiles) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *PostFiles) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

// CirclePostModel 圈子帖子表
type CirclePostModel struct {
	models.Model
	PostID    string     `gorm:"column:post_id;size:64;uniqueIndex;not null" json:"postId"` // 帖子唯一ID
	CircleID  string     `gorm:"size:64;not null;index" json:"circleId"`                    // 所属圈子ID
	UserID    string     `gorm:"size:64;not null;index" json:"userId"`                      // 发帖用户ID
	Content   string     `gorm:"type:text;not null" json:"content"`                         // 帖子内容
	Files     *PostFiles `gorm:"type:longtext" json:"files"`                                // 附件列表（JSON数组）
	IsTop     bool       `gorm:"not null;default:false" json:"isTop"`                       // 是否置顶
	IsDeleted bool       `gorm:"not null;default:false;index" json:"isDeleted"`             // 软删除标记
}
