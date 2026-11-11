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

package open_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// DefaultKnowledgeSourceTopK 知识源默认召回条数。
const DefaultKnowledgeSourceTopK = 8

// MaxKnowledgeSourceTopK 召回条数上限，避免一次拉回过多片段。
const MaxKnowledgeSourceTopK = 50

// OpenKnowledgeSource 知识源：智能体侧的一等实体，归某个数字伙伴所有。
//
// 知识源只持有「引用」，不存任何文档正文：
//   - 文档正文的权威源在 doc 域在线文档库（普通用户所有，按 user_id 鉴权）；
//   - 新建知识源时由 portal 在 platform 域建一个专属索引集合（kind=open_knowledge_source），
//     用户粘贴的在线文档地址经解析 + 鉴权后，作为「引用」挂在这个集合上，
//     由 platform 的 ingest 自动拉取并索引源文档；
//   - 删除知识源 / 移除引用都不会动在线文档本体（飞书/钉钉「引用云文档」语义）。
//
// KnowledgeIDs 只放 1 个 platform 索引集合 ID，不存任何文档正文。
// 运行时不关心分组：BuiltinConfig.KnowledgeIDs 由本表所有行合并去重后派生，
// 知识源发生任何增删改都由 portal logic 回写，agent 侧只认那个扁平列表。
// 控制台的 config/save 会整列覆盖 builtin_config，因此也在那里按本表重算回填。
type OpenKnowledgeSource struct {
	gorm.Model
	SourceID      string          `gorm:"column:source_id;type:varchar(64);uniqueIndex;not null;comment:知识源ID（ks_ 前缀）"`
	DigitalID     string          `gorm:"column:digital_id;type:varchar(64);not null;index;comment:所属数字伙伴ID（OpenDigital.DigitalID）"`
	OwnerID       string          `gorm:"column:owner_id;type:varchar(64);not null;index;comment:开发者用户ID"`
	Name          string          `gorm:"column:name;type:varchar(128);not null;comment:知识源名称"`
	RetrievalTopK int             `gorm:"column:retrieval_top_k;type:int;not null;default:8;comment:召回条数"`
	IndexConfig   string          `gorm:"column:index_config;type:text;comment:分段与清洗配置 JSON"`
	KnowledgeIDs  KnowledgeIDList `gorm:"column:knowledge_ids;type:json;comment:持有的 platform 索引集合 ID（1 个）"`
}

func (OpenKnowledgeSource) TableName() string { return "open_knowledge_source" }

// KnowledgeIDList 是 JSON 列里的字符串数组。
// 单独定义 Scan/Value 而不是依赖 serializer tag：库里可能是 NULL、空串或 JSON 数组，
// 三种情况都要稳定读成空切片，避免 logic 里到处判 nil。
type KnowledgeIDList []string

func (l KnowledgeIDList) Value() (driver.Value, error) {
	if len(l) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(l)
}

func (l *KnowledgeIDList) Scan(value interface{}) error {
	if value == nil {
		*l = KnowledgeIDList{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}
	if len(bytes) == 0 {
		*l = KnowledgeIDList{}
		return nil
	}
	return json.Unmarshal(bytes, l)
}
