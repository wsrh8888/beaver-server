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
 *
 * 场景知识源：场景智能体侧的一等实体，归某个固定场景（assistant|chat|moment|workbench）所有。
 *
 * 与开放平台的 open_knowledge_source 同构，差别只在归属主体：
 *   - open 端归属数字伙伴（digital_id），因为那是开发者自己创建的应用；
 *   - 这里归属平台固定场景（scene），场景由平台内置，运营在后台配置。
 *
 * 知识源只持有「引用」，不存任何文档正文：
 *   - 新建知识源时在 platform 域建一个专属索引集合，本表只存这个集合 ID；
 *   - 运营粘贴的在线文档地址经解析后，作为 source_type=doc 的引用挂在该集合上，
 *     由 platform 的 ingest 自动拉取并索引源文档（引用云文档语义）；
 *   - 本地上传的文件解析成正文后以 source_type=file 写入同一个集合。
 *
 * 运行时只认 agent_apps.knowledge_ids：本表任何增删改都由后端把各知识源的索引集合
 * 合并去重后回写，agent 侧不感知分组。
 */

package backend_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"beaver/common/models"
)

// DefaultAgentKnowledgeSourceTopK 场景知识源默认召回条数。
const DefaultAgentKnowledgeSourceTopK = 8

// MaxAgentKnowledgeSourceTopK 召回条数上限，避免一次拉回过多片段。
const MaxAgentKnowledgeSourceTopK = 50

// AgentKnowledgeSource 场景知识源实体。表 agent_knowledge_source。
type AgentKnowledgeSource struct {
	models.Model
	SourceID      string          `gorm:"column:source_id;size:64;uniqueIndex;not null;comment:知识源ID（ks_ 前缀）"`
	Scene         string          `gorm:"column:scene;size:32;not null;index;comment:所属场景（assistant|chat|moment|workbench）"`
	Name          string          `gorm:"column:name;size:128;not null;comment:知识源名称"`
	RetrievalTopK int             `gorm:"column:retrieval_top_k;not null;default:8;comment:召回条数"`
	IndexConfig   string          `gorm:"column:index_config;type:text;comment:分段与清洗配置 JSON"`
	KnowledgeIDs  KnowledgeIDList `gorm:"column:knowledge_ids;type:json;comment:持有的 platform 索引集合 ID（1 个）"`
}

func (AgentKnowledgeSource) TableName() string { return "agent_knowledge_source" }

// KnowledgeIDList 是 JSON 列里的字符串数组。
// 单独定义 Scan/Value：库里可能是 NULL、空串或 JSON 数组，三种情况都要稳定读成空切片。
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
