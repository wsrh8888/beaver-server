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

package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/go-redis/redis"

	"gorm.io/gorm"
)

// VersionGenerator 版本号生成器（Redis缓存 + MySQL主存储）
type VersionGenerator struct {
	redisClient *redis.Client
	db          *gorm.DB
}

// NewVersionGenerator 创建版本号生成器
func NewVersionGenerator(redisClient *redis.Client, db *gorm.DB) *VersionGenerator {
	return &VersionGenerator{
		redisClient: redisClient,
		db:          db,
	}
}

// GetNextVersion 获取下一个版本号（Redis缓存 + MySQL主存储）
// 参数:
//   - table: 表名或业务标识符 (必传)
//   - field: 条件字段名；为空表示全局版本。若包含下划线，表示多字段组合键（如 user_id_category），会按下划线拆分为多列。
//   - value: 条件字段值；为空表示全局版本。若 field 为组合键，value 也需用下划线分段（如 uid_cat），与 field 拆分后的段数一致。
//
// 返回: 版本号，失败时返回-1
func (vg *VersionGenerator) GetNextVersion(table string, field string, value string) int64 {
	var key string
	isConditional := field != "" && value != ""

	if isConditional {
		key = fmt.Sprintf("version:%s:%s:%s", table, field, value)
	} else {
		key = fmt.Sprintf("version:%s", table)
	}

	// 默认从1开始
	initValue := int64(1)

	// 1. 先尝试从Redis获取
	version, err := vg.redisClient.Incr(key).Result()
	if err != nil {
		beaverlog.New("version").Error(model.LogMsg{
			Text: "Redis获取版本号失败",
			Data: map[string]any{"table": table, "field": field, "value": value, "err": err.Error()},
		})
		// Redis失败，初始化Redis并返回初始值
		vg.redisClient.Set(key, initValue, 7*24*time.Hour) // 7天过期
		beaverlog.New("version").Info(model.LogMsg{
			Text: "初始化版本号",
			Data: map[string]any{"table": table, "field": field, "value": value, "initValue": initValue},
		})
		return initValue
	}

	// 2. 检查Redis版本号是否合理（防止Redis重启或过期后从0开始）
	if version == 1 {
		// Redis返回1表示key不存在或刚创建，需要从MySQL同步历史版本
		mysqlVersion := vg.getMaxVersionFromMySQL(table, field, value)
		if mysqlVersion < 0 {
			beaverlog.New("version").Error(model.LogMsg{
				Text: "MySQL获取版本号失败",
				Data: map[string]any{"table": table},
			})
			// MySQL查询失败，继续使用Redis的版本号1
			return version
		}

		if mysqlVersion > 0 {
			// MySQL有数据，更新Redis为MySQL版本号+1
			newVersion := mysqlVersion + 1
			vg.redisClient.Set(key, newVersion, 7*24*time.Hour) // 7天过期
			beaverlog.New("version").Info(model.LogMsg{
				Text: "从MySQL同步版本号",
				Data: map[string]any{"table": table, "mysqlVersion": mysqlVersion, "newVersion": newVersion},
			})
			return newVersion
		}

		// MySQL没有数据，使用初始值
		if initValue > 1 {
			vg.redisClient.Set(key, initValue, 7*24*time.Hour) // 7天过期
			beaverlog.New("version").Info(model.LogMsg{
				Text: "使用自定义初始值",
				Data: map[string]any{"table": table, "initValue": initValue},
			})
			return initValue
		}
		// MySQL没有数据，继续使用Redis的版本号1
		beaverlog.New("version").Info(model.LogMsg{
			Text: "MySQL没有数据Redis从1开始",
			Data: map[string]any{"table": table, "version": version},
		})
	}

	beaverlog.New("version").Info(model.LogMsg{
		Text: "生成版本号",
		Data: map[string]any{"table": table, "field": field, "value": value, "version": version},
	})
	return version
}

// getMaxVersionFromMySQL 从MySQL获取最大版本号
// 对于条件版本，查询该条件下的最大版本号，失败时返回-1
func (vg *VersionGenerator) getMaxVersionFromMySQL(table string, field string, value string) int64 {
	isConditional := field != "" && value != ""
	tableName := vg.getTableName(table)

	var maxVersion int64
	var err error

	if isConditional {
		if strings.Contains(field, "_") && strings.Contains(value, "_") {
			fields := strings.Split(field, "_")
			values := strings.Split(value, "_")
			if len(fields) != len(values) {
				beaverlog.New("version").Error(model.LogMsg{
					Text: "组合键字段和值数量不匹配",
					Data: map[string]any{"fields": fields, "values": values},
				})
				return -1
			}
			conds := make([]string, 0, len(fields))
			args := make([]interface{}, 0, len(fields))
			for i := range fields {
				if !isSafeIdentifier(fields[i]) {
					beaverlog.New("version").Error(model.LogMsg{
						Text: "非法字段名",
						Data: map[string]any{"field": fields[i]},
					})
					return -1
				}
				conds = append(conds, fmt.Sprintf("%s = ?", fields[i]))
				args = append(args, values[i])
			}
			query := fmt.Sprintf("SELECT COALESCE(MAX(version), 0) FROM %s WHERE %s", tableName, strings.Join(conds, " AND "))
			err = vg.db.Raw(query, args...).Scan(&maxVersion).Error
			beaverlog.New("version").Info(model.LogMsg{
				Text: "查询组合条件版本历史",
				Data: map[string]any{"table": table, "fields": fields, "values": values, "maxVersion": maxVersion},
			})
		} else {
			// 单字段条件
			if !isSafeIdentifier(field) {
				beaverlog.New("version").Error(model.LogMsg{
					Text: "非法字段名",
					Data: map[string]any{"field": field},
				})
				return -1
			}
			query := fmt.Sprintf("SELECT COALESCE(MAX(version), 0) FROM %s WHERE %s = ?", tableName, field)
			err = vg.db.Raw(query, value).Scan(&maxVersion).Error
			beaverlog.New("version").Info(model.LogMsg{
				Text: "查询条件版本历史",
				Data: map[string]any{"table": table, "field": field, "value": value, "maxVersion": maxVersion},
			})
		}
	} else {
		// 全局版本：查询表中的最大版本号
		err = vg.db.Raw("SELECT COALESCE(MAX(version), 0) FROM " + tableName).Scan(&maxVersion).Error
		beaverlog.New("version").Info(model.LogMsg{
			Text: "查询全局版本历史",
			Data: map[string]any{"table": table, "maxVersion": maxVersion},
		})
	}

	if err != nil {
		beaverlog.New("version").Error(model.LogMsg{
			Text: "MySQL获取最大版本号失败",
			Data: map[string]any{"table": table, "field": field, "value": value, "err": err.Error()},
		})
		return -1
	}
	return maxVersion
}

// getTableName 根据数据类型获取表名
// 注意：条件版本不应该调用此方法，因为条件版本不从MySQL查询
func (vg *VersionGenerator) getTableName(dataType string) string {
	switch dataType {
	case "users":
		return "user_models"
	case "friends":
		return "friend_models"
	case "groups":
		return "group_models"
	case "circles":
		return "circle_models"
	case "group_members":
		return "group_member_models"
	case "group_member_logs":
		return "group_member_change_log_models"
	case "chats":
		return "message_models"
	case "notification_events":
		return "notification_events"
	case "notification_inboxes":
		return "notification_inboxes"
	case "notification_reads":
		return "notification_reads"
	default:
		return dataType
	}
}

var safeIdentRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func isSafeIdentifier(s string) bool {
	return safeIdentRegexp.MatchString(s)
}
