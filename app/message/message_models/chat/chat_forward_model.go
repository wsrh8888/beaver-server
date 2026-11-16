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

package chat_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ForwardContent 转发内容集合类型，实现 GORM 的 Scan 和 Value 接口
type ForwardContent []ChatMessage

func (m ForwardContent) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *ForwardContent) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

// ChatForwardDetail 合并转发详情表（冷热分离，存储大块JSON数据）
type ChatForward struct {
	models.Model
	RecordID string         `gorm:"size:64;uniqueIndex" json:"recordId"` // 聚合ID
	Content  ForwardContent `gorm:"type:json" json:"content"`            // 序列化后的消息数组快照
}
