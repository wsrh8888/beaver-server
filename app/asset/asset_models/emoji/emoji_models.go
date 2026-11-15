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

package emoji_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// EmojiInfo 表情信息
type EmojiInfo struct {
	Width  int `json:"width"`  // 表情图片宽度
	Height int `json:"height"` // 表情图片高度
}

// Value converts the EmojiInfo to a JSON-encoded string for database storage
func (e *EmojiInfo) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Scan converts a JSON-encoded string from the database to a EmojiInfo
func (e *EmojiInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, e)
}

// 表情
type Emoji struct {
	models.Model
	EmojiID   string    `gorm:"column:emoji_id;size:64;uniqueIndex" json:"emojiId"` // 全局唯一ID
	FileKey   string    `json:"fileKey"`                                            // 文件Key
	Title     string    `json:"title"`                                              // 表情名称
	EmojiInfo EmojiInfo `gorm:"type:longtext;serializer:json" json:"emojiInfo"`    // 表情详细信息（JSON格式）
	Status    int8      `gorm:"default:1" json:"status"`                            // 状态：1=正常 2=审核中 3=违规禁用
	Version   int64     `gorm:"not null;default:0;index" json:"version"`            //基于表递增
}
