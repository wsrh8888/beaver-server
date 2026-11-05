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

import "beaver/common/models"

// KnowledgeChunk 文档分段。表 knowledge_chunks。
// 向量本体存 OpenSearch，DB 只存正文 + 段号，对齐 stratum 的「文档+段号」可回溯检索：
// 命中分片靠 document_id + chunk_index 回到原文定位，embedding_ref 记录向量在索引中的落点。
type KnowledgeChunk struct {
	models.Model
	ChunkID      string `gorm:"column:chunk_id;type:varchar(64);uniqueIndex;not null;comment:分段业务ID" json:"chunkId"`
	CollectionID string `gorm:"column:collection_id;type:varchar(64);not null;index;comment:所属集合ID（冗余便于按库清理）" json:"collectionId"`
	DocumentID   string `gorm:"column:document_id;type:varchar(64);not null;index;comment:所属文档ID" json:"documentId"`
	ChunkIndex   int    `gorm:"column:chunk_index;type:int;not null;default:0;comment:段序号，文档内从0递增" json:"chunkIndex"`
	Content      string `gorm:"type:longtext;comment:分段正文" json:"content"`
	CharCount    int    `gorm:"column:char_count;type:int;not null;default:0;comment:正文字符数" json:"charCount"`
	EmbeddingRef string `gorm:"column:embedding_ref;type:varchar(255);not null;default:'';comment:向量在 OpenSearch 的落点（index#docId）" json:"embeddingRef"`
}

func (KnowledgeChunk) TableName() string { return "knowledge_chunks" }
