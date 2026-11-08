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

package knowledge

import (
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

// toCollectionInfo 把 platform_rpc 集合项映射为后台 API 结构体。
func toCollectionInfo(c *platform_rpc.KnowledgeCollectionItem) types.KnowledgeCollectionInfo {
	if c == nil {
		return types.KnowledgeCollectionInfo{}
	}
	return types.KnowledgeCollectionInfo{
		Id:           int64(c.Id),
		CollectionID: c.CollectionId,
		ParentID:     c.ParentId,
		Name:         c.Name,
		Kind:         c.Kind,
		AclMode:      c.AclMode,
		IndexConfig:  c.IndexConfig,
		Description:  c.Description,
		OwnerUserID:  c.OwnerUserId,
		Status:       int(c.Status),
		Sort:         int(c.Sort),
	}
}

// toDocumentInfo 把 platform_rpc 文档项映射为后台 API 结构体（不含正文）。
func toDocumentInfo(d *platform_rpc.KnowledgeDocumentItem) types.KnowledgeDocumentInfo {
	if d == nil {
		return types.KnowledgeDocumentInfo{}
	}
	return types.KnowledgeDocumentInfo{
		Id:              int64(d.Id),
		DocumentID:      d.DocumentId,
		CollectionID:    d.CollectionId,
		Title:           d.Title,
		MediaType:       d.MediaType,
		DocType:         d.DocType,
		ByteSize:        d.ByteSize,
		SourceType:      d.SourceType,
		SourceDocId:     d.SourceDocId,
		SourceRevision:  d.SourceRevision,
		SyncMode:        d.SyncMode,
		SyncInterval:    int(d.SyncInterval),
		LastSyncedAt:    d.LastSyncedAt,
		MetadataText:    d.MetadataText,
		RagEnabled:      d.RagEnabled,
		IndexStatus:     d.IndexStatus,
		IndexRetryCount: int(d.IndexRetryCount),
		IndexError:      d.IndexError,
		IndexedAt:       d.IndexedAt,
		Status:          int(d.Status),
	}
}
