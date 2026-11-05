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
	"time"

	"beaver/common/models"
)

// KnowledgeDocument 知识文档。表 knowledge_documents。
// 索引状态机：pending → indexing → ready / failed；worker 按 index_next_at 捞取、抢占式 UPDATE 去重、
// 超上限置 failed 并留 index_error。删除走软删（DeletedAt），删时同步清 OpenSearch 索引。
// 为将来「在线多人文档」边界预留：owner_user_id / version / updated_by（本期不实现 OT/CRDT）。
type KnowledgeDocument struct {
	models.Model
	DocumentID   string `gorm:"column:document_id;type:varchar(64);uniqueIndex;not null;comment:文档业务ID" json:"documentId"`
	CollectionID string `gorm:"column:collection_id;type:varchar(64);not null;index;comment:所属集合ID" json:"collectionId"`
	Title        string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	MediaType    string `gorm:"column:media_type;type:varchar(64);not null;default:'text/plain';comment:MIME 类型" json:"mediaType"`
	DocType      string `gorm:"column:doc_type;type:varchar(32);not null;default:'guide';index;comment:guide|policy|faq|product" json:"docType"`
	ByteSize     int64  `gorm:"column:byte_size;type:bigint;not null;default:0;comment:原始字节数" json:"byteSize"`
	SourceRef    string `gorm:"column:source_ref;type:varchar(512);not null;default:'';comment:对象存储键" json:"sourceRef"`
	ContentText  string `gorm:"column:content_text;type:longtext;comment:解析出的正文（分段前原文）" json:"contentText"`
	ContentDigest string `gorm:"column:content_digest;type:varchar(64);not null;default:'';comment:正文哈希，判重/免重复索引" json:"contentDigest"`
	MetadataText string `gorm:"column:metadata_text;type:json;comment:文档级元数据 JSON" json:"metadataText"`
	RagEnabled   bool   `gorm:"column:rag_enabled;type:tinyint;not null;default:1;comment:是否参与检索" json:"ragEnabled"`
	OwnerUserID  string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户（多人文档预留）" json:"ownerUserId"`
	Version      int    `gorm:"type:int;not null;default:1;comment:版本号（多人文档预留）" json:"version"`
	UpdatedBy    string `gorm:"column:updated_by;type:varchar(64);not null;default:'';comment:最后编辑者（多人文档预留）" json:"updatedBy"`

	// —— 索引状态机 ——
	IndexStatus    string     `gorm:"column:index_status;type:varchar(16);not null;default:'pending';index;comment:pending|indexing|ready|failed" json:"indexStatus"`
	IndexRetryCount int       `gorm:"column:index_retry_count;type:int;not null;default:0;comment:已重试次数" json:"indexRetryCount"`
	IndexNextAt    *time.Time `gorm:"column:index_next_at;type:datetime;index;comment:下次可索引时间（退避调度）" json:"indexNextAt"`
	IndexError     string     `gorm:"column:index_error;type:varchar(512);not null;default:'';comment:最近一次失败原因" json:"indexError"`
	IndexedAt      *time.Time `gorm:"column:indexed_at;type:datetime;comment:索引完成时间" json:"indexedAt"`
	Status         int        `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用（停用不参与检索）" json:"status"`

	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;not null;default:'1970-01-01 00:00:00';index;comment:软删时间，1970 为未删" json:"deletedAt"`
}

func (KnowledgeDocument) TableName() string { return "knowledge_documents" }

// IsDeleted 软删判定：deleted_at 非零值即视为已删。
func (d *KnowledgeDocument) IsDeleted() bool {
	return !d.DeletedAt.IsZero() && d.DeletedAt.Year() > 1970
}
