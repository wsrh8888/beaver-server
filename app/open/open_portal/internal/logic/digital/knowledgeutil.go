/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

// loadDigitalKnowledgeIDs 读数字伙伴运行时检索用的知识库 ID 扁平列表。
// 知识库本体由 platform 域维护，open 侧只存 ID（跨服务不复制内容）。
//
// 这列 ID 的唯一来源是 open_knowledge_source 表：每个知识源可引入 N 个知识库，
// 由 loadKnowledgeSourceIDs 合并去重后写入。agent 侧不关心知识源分组，只认这个列表。
func loadDigitalKnowledgeIDs(db *gorm.DB, digitalID string) []string {
	var cfg open_models.OpenDigitalConfig
	if err := db.Where("digital_id = ?", digitalID).First(&cfg).Error; err != nil {
		return nil
	}
	if cfg.BuiltinConfig == nil {
		return nil
	}
	return cfg.BuiltinConfig.KnowledgeIDs
}

// writeDigitalKnowledgeIDs 覆盖写入数字伙伴的绑定知识库列表。
// 没有配置行时先建空行：新建的数字伙伴可能还没保存过编排配置。
func writeDigitalKnowledgeIDs(db *gorm.DB, digitalID string, ids []string) error {
	return withBuiltinConfig(db, digitalID, func(cfg *open_models.BuiltinConfig) {
		cfg.KnowledgeIDs = ids
	})
}

// loadDigitalMcpServerIDs 读数字伙伴启用的 MCP 服务 ID。
func loadDigitalMcpServerIDs(db *gorm.DB, digitalID string) []string {
	var cfg open_models.OpenDigitalConfig
	if err := db.Where("digital_id = ?", digitalID).First(&cfg).Error; err != nil {
		return nil
	}
	if cfg.BuiltinConfig == nil {
		return nil
	}
	return cfg.BuiltinConfig.McpServerIDs
}

// writeDigitalMcpServerIDs 覆盖写入启用的 MCP 服务 ID。
func writeDigitalMcpServerIDs(db *gorm.DB, digitalID string, ids []string) error {
	return withBuiltinConfig(db, digitalID, func(cfg *open_models.BuiltinConfig) {
		cfg.McpServerIDs = ids
	})
}

// withBuiltinConfig 只改 builtin 配置的一小部分，其余字段（人设、工具、知识…）原样保留。
func withBuiltinConfig(db *gorm.DB, digitalID string, mutate func(*open_models.BuiltinConfig)) error {
	var cfg open_models.OpenDigitalConfig
	err := db.Where("digital_id = ?", digitalID).First(&cfg).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg = open_models.OpenDigitalConfig{DigitalID: digitalID}
	}
	if cfg.BuiltinConfig == nil {
		cfg.BuiltinConfig = &open_models.BuiltinConfig{}
	}
	mutate(cfg.BuiltinConfig)
	return db.Save(&cfg).Error
}

// indexConfigTopK 从 index_config JSON 取检索条数，缺省返回 0（= 平台默认）。
func indexConfigTopK(raw string) int {
	m := parseJSONMap(raw)
	for _, key := range []string{"top_k", "topK"} {
		switch v := m[key].(type) {
		case float64:
			return int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return 0
}

// mergeIndexConfig 把改动合并进原有 index_config。
// 覆盖式写入会把 embedding_endpoint / model 等凭据抹掉，所以只覆盖传入的键。
func mergeIndexConfig(existing, patch string) string {
	base := parseJSONMap(existing)
	overlay := parseJSONMap(patch)
	if len(overlay) == 0 {
		return existing
	}
	for k, v := range overlay {
		base[k] = v
	}
	merged, err := json.Marshal(base)
	if err != nil {
		return existing
	}
	return string(merged)
}

func parseJSONMap(raw string) map[string]interface{} {
	out := map[string]interface{}{}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

// setIndexConfigTopK 在 index_config 上写 top_k（设置页只改检索条数时用）。
func setIndexConfigTopK(raw string, topK int) string {
	if topK <= 0 {
		return raw
	}
	return mergeIndexConfig(raw, `{"top_k":`+strconv.Itoa(topK)+`}`)
}

// parseRFC3339Unix 把 RPC 返回的时间字符串转成秒级时间戳，解析失败给 0。
func parseRFC3339Unix(raw string) int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Unix()
		}
	}
	return 0
}
