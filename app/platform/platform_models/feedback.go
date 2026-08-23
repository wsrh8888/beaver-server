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

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// FeedbackType 反馈类型
type FeedbackType int

const (
	FeedbackTypeBug     FeedbackType = 1 // 问题反馈
	FeedbackTypeFeature FeedbackType = 2 // 功能建议
	FeedbackTypeOther   FeedbackType = 3 // 其他
)

// FeedbackStatus 反馈状态
type FeedbackStatus int

const (
	FeedbackStatusPending    FeedbackStatus = 1 // 待处理
	FeedbackStatusProcessing FeedbackStatus = 2 // 处理中
	FeedbackStatusResolved   FeedbackStatus = 3 // 已解决
	FeedbackStatusRejected   FeedbackStatus = 4 // 已拒绝
)

// FileNames 文件ID数组类型
type FileNames []string

// Value 实现 driver.Valuer 接口，用于写入数据库
func (f FileNames) Value() (driver.Value, error) {
	if f == nil {
		return "[]", nil
	}
	return json.Marshal(f)
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据
func (f *FileNames) Scan(value interface{}) error {
	if value == nil {
		*f = FileNames{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

// FeedbackModel 反馈表
type FeedbackModel struct {
	models.Model
	UserID       string         `gorm:"index" json:"user_id"`                          // 提交用户ID
	Content      string         `gorm:"type:text;not null" json:"content"`             // 反馈内容
	Type         FeedbackType   `gorm:"type:tinyint;not null" json:"type"`             // 反馈类型
	Status       FeedbackStatus `gorm:"type:tinyint;not null;default:1" json:"status"` // 反馈状态
	FileNames    FileNames      `gorm:"type:json" json:"fileNames"`                    // 反馈图片URL数组
	HandlerID    int64          `gorm:"index" json:"handlerId"`                        // 处理人ID
	HandleTime   *time.Time     `json:"handle_time"`                                   // 处理时间
	HandleResult string         `gorm:"type:text" json:"handleResult"`                 // 处理结果
}
