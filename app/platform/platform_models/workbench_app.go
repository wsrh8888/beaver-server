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
)

// WorkbenchEntryConfig 入口配置（JSON 落库）
// type=0：pc/mobile 填路由 key，如 moment
// type=1：pc/mobile 填 H5 URL
type WorkbenchEntryConfig struct {
	Type   int8   `json:"type"`             // 0路由 1URL
	PC     string `json:"pc,omitempty"`     // PC 入口
	Mobile string `json:"mobile,omitempty"` // 移动端入口
}

func (c WorkbenchEntryConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *WorkbenchEntryConfig) Scan(value interface{}) error {
	if value == nil {
		*c = WorkbenchEntryConfig{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

// WorkbenchApp IM 运营工作台应用（后台配置、客户端展示）
type WorkbenchApp struct {
	models.Model
	WorkbenchAppID string               `gorm:"size:64;uniqueIndex;not null;comment:应用业务ID" json:"workbenchAppId"`
	Name           string               `gorm:"size:64;not null;comment:应用名称" json:"name"`
	Description    string               `gorm:"size:500;comment:应用描述" json:"description"`
	Icon           string               `gorm:"size:500;comment:应用图标URL" json:"icon"`
	AppType        int8                 `gorm:"type:tinyint;not null;default:1;index;comment:类型 0内部 1第三方H5" json:"appType"`
	ClientScope    int8                 `gorm:"type:tinyint;not null;default:0;index;comment:可见端 0全部 1仅PC 2仅移动" json:"clientScope"`
	EntryConfig    WorkbenchEntryConfig `gorm:"type:json;not null;comment:入口配置JSON" json:"entryConfig"`
	OpenMode       int8                 `gorm:"type:tinyint;not null;default:0;comment:打开方式 0内嵌 1系统浏览器(仅第三方)" json:"openMode"`
	Category       int8                 `gorm:"type:tinyint;not null;default:0;index;comment:分组 0默认 1办公 2审批 3效率 4其他" json:"category"`
	Sort           int                  `gorm:"not null;default:0;index;comment:排序(越小越靠前)" json:"sort"`
	Status         int8                 `gorm:"type:tinyint;not null;default:0;index;comment:状态 0下架 1上架" json:"status"`
	CreatedBy      string               `gorm:"size:64;comment:创建人(管理员ID)" json:"createdBy"`
	LastModifiedBy string               `gorm:"size:64;comment:最后修改人(管理员ID)" json:"lastModifiedBy"`
	Remark         string               `gorm:"size:500;comment:运营备注(不对客户端暴露)" json:"remark"`
}

func (WorkbenchApp) TableName() string {
	return "workbench_apps"
}
