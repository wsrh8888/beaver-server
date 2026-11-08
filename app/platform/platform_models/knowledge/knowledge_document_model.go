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

package knowledge_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"beaver/common/models"
)

// KnowledgeSourceDetail 知识文档来源明细，按 JSON 存一列，减少零散字段。
// 这些内容只用于导入/自动同步，不参与检索，所以不单独建列。
type KnowledgeSourceDetail struct {
	DocID        string     `json:"docId,omitempty"`        // 平台在线文档ID（doc 业务的 doc_id）
	Revision     string     `json:"revision,omitempty"`     // 来源版本（doc version）
	SyncInterval int        `json:"syncInterval,omitempty"` // 自动同步间隔秒
	LastSyncedAt *time.Time `json:"lastSyncedAt,omitempty"` // 最近同步时间
}

// Value 把来源明细序列化为 JSON 存库。
func (s KnowledgeSourceDetail) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 从库里的 JSON 反序列化来源明细。
func (s *KnowledgeSourceDetail) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// KnowledgeDocument 知识文档。表 knowledge_documents。
// 索引状态机：pending → indexing → ready / failed；worker 按 index_next_at 捞取、抢占式 UPDATE 去重、
// 超上限置 failed 并留 index_error。删除走软删（DeletedAt），删时同步清 OpenSearch 索引。
// 本表只保存导入后的知识快照，不参与在线编辑。来源只有两种：doc（本平台在线文档）、file（上传文件）。
// 自动同步模式下按 source_detail.docId + doc version 增量刷新，权限跟随集合 acl_mode。
type KnowledgeDocument struct {
	models.Model
	DocumentID    string                `gorm:"column:document_id;type:varchar(64);uniqueIndex;not null;comment:文档业务ID" json:"documentId"`
	CollectionID  string                `gorm:"column:collection_id;type:varchar(64);not null;index;comment:所属集合ID" json:"collectionId"`
	Title         string                `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	MediaType     string                `gorm:"column:media_type;type:varchar(64);not null;default:'text/plain';comment:MIME 类型" json:"mediaType"`
	DocType       string                `gorm:"column:doc_type;type:varchar(32);not null;default:'guide';index;comment:guide|policy|faq|product" json:"docType"`
	ByteSize      int64                 `gorm:"column:byte_size;type:bigint;not null;default:0;comment:原始字节数" json:"byteSize"`
	SourceType    string                `gorm:"column:source_type;type:varchar(16);not null;default:'file';index;comment:doc平台在线文档|file上传文件" json:"sourceType"`
	SyncMode      string                `gorm:"column:sync_mode;type:varchar(16);not null;default:'import';index;comment:import单次导入|direct自动同步" json:"syncMode"`
	SourceDetail  KnowledgeSourceDetail `gorm:"column:source_detail;type:longtext;comment:来源明细 JSON（文档ID/版本/同步间隔/最近同步时间）" json:"sourceDetail"`
	ContentText   string                `gorm:"column:content_text;type:longtext;comment:解析出的正文（分段前原文）" json:"contentText"`
	ContentDigest string                `gorm:"column:content_digest;type:varchar(64);not null;default:'';comment:正文哈希，判重/免重复索引" json:"contentDigest"`
	MetadataText  string                `gorm:"column:metadata_text;type:json;comment:文档级元数据 JSON" json:"metadataText"`
	RagEnabled    bool                  `gorm:"column:rag_enabled;type:tinyint;not null;default:1;comment:是否参与检索" json:"ragEnabled"`

	// —— 索引状态机 ——
	IndexStatus     string     `gorm:"column:index_status;type:varchar(16);not null;default:'pending';index;comment:pending|indexing|ready|failed" json:"indexStatus"`
	IndexRetryCount int        `gorm:"column:index_retry_count;type:int;not null;default:0;comment:已重试次数" json:"indexRetryCount"`
	IndexNextAt     *time.Time `gorm:"column:index_next_at;type:datetime;index;comment:下次可索引时间（退避调度）" json:"indexNextAt"`
	IndexError      string     `gorm:"column:index_error;type:varchar(512);not null;default:'';comment:最近一次失败原因" json:"indexError"`
	IndexedAt       *time.Time `gorm:"column:indexed_at;type:datetime;comment:索引完成时间" json:"indexedAt"`
	Status          int        `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用（停用不参与检索）" json:"status"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间，1970 为未删" json:"deletedAt"`
}

func (KnowledgeDocument) TableName() string { return "knowledge_documents" }

// IsDeleted 软删判定：deleted_at 非零值即视为已删。
func (d *KnowledgeDocument) IsDeleted() bool {
	return !d.DeletedAt.IsZero() && d.DeletedAt.Year() > 1970
}
