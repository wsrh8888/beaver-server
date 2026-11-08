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

package worker

import (
	"context"
	"time"

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/ingest"
	"beaver/app/platform/platform_rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// resyncOnce 处理自动同步（direct）的在线文档：到期的按 doc_id 重新读取，正文变化才重建索引。
func resyncOnce(ctx context.Context, svcCtx *svc.ServiceContext) {
	batch := svcCtx.Config.Index.BatchSize
	if batch <= 0 {
		batch = 16
	}
	var docs []knowledge_models.KnowledgeDocument
	err := svcCtx.DB.
		Where("sync_mode = ? AND source_type = ? AND status = ? AND deleted_at = ?", "direct", "doc", 1, zeroDeletedAt).
		Order("updated_at ASC").Limit(batch).Find(&docs).Error
	if err != nil {
		logx.Errorf("捞取待同步知识文档失败: %v", err)
		return
	}
	for i := range docs {
		syncDocument(ctx, svcCtx, &docs[i])
	}
}

func syncDocument(ctx context.Context, svcCtx *svc.ServiceContext, doc *knowledge_models.KnowledgeDocument) {
	now := time.Now()
	interval := time.Duration(doc.SourceDetail.SyncInterval) * time.Second
	if interval <= 0 || doc.SourceDetail.DocID == "" {
		return
	}
	if doc.SourceDetail.LastSyncedAt != nil && now.Before(doc.SourceDetail.LastSyncedAt.Add(interval)) {
		return
	}
	res, err := ingest.Fetch(ctx, svcCtx.DocRpc, ingest.Source{Type: doc.SourceType, DocID: doc.SourceDetail.DocID})
	if err != nil {
		logx.Errorf("自动同步知识文档失败 document_id=%s: %v", doc.DocumentID, err)
		return
	}
	src := doc.SourceDetail
	src.Revision = res.Revision
	src.LastSyncedAt = &now
	updates := map[string]any{
		"source_detail": src,
	}
	digest := ingest.Digest(res.Content)
	if digest != doc.ContentDigest {
		updates["content_text"] = res.Content
		updates["content_digest"] = digest
		updates["byte_size"] = int64(len([]byte(res.Content)))
		updates["index_status"] = "pending"
		updates["index_next_at"] = nil
		updates["index_error"] = ""
	}
	if err := svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ?", doc.DocumentID).Updates(updates).Error; err != nil {
		logx.Errorf("回写知识文档同步状态失败 document_id=%s: %v", doc.DocumentID, err)
	}
}
