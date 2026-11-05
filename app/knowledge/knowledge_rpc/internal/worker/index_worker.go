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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/embedding"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/core/coreopensearch"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// zeroDeletedAt 软删哨兵：deleted_at 等于该时刻表示未删，worker 只捞未删文档。
const zeroDeletedAt = "1970-01-01 00:00:00"

// Start 拉起进程内索引 worker：先确保 OpenSearch 索引 mapping 存在，再按周期轮询待索引文档。
// 照搬 stratum：pending 捞取 + 抢占式 UPDATE 去重 + 重试上限置 failed。ctx 取消即退出。
func Start(ctx context.Context, svcCtx *svc.ServiceContext) {
	if !svcCtx.Config.Index.Enabled {
		logx.Info("知识索引 worker 已禁用")
		return
	}
	if svcCtx.OpenSearch == nil {
		logx.Error("知识索引 worker 未启动：OpenSearch 未配置")
		return
	}
	dims := svcCtx.Config.Index.Dims
	if dims <= 0 {
		dims = 1536
	}
	if err := svcCtx.OpenSearch.EnsureKnowledgeIndex(ctx, svcCtx.ChunkIndex, dims); err != nil {
		logx.Errorf("确保知识索引 mapping 失败（worker 仍启动，检索期再重试）: %v", err)
	}
	interval := svcCtx.Config.Index.IntervalSeconds
	if interval <= 0 {
		interval = 5
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	logx.Infof("知识索引 worker 已启动，轮询间隔 %ds，索引 %s", interval, svcCtx.ChunkIndex)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runOnce(ctx, svcCtx)
		}
	}
}

// runOnce 捞取一批 pending 文档并逐个抢占处理。
func runOnce(ctx context.Context, svcCtx *svc.ServiceContext) {
	batch := svcCtx.Config.Index.BatchSize
	if batch <= 0 {
		batch = 16
	}
	var docs []knowledge_models.KnowledgeDocument
	now := time.Now()
	err := svcCtx.DB.
		Where("index_status = ? AND deleted_at = ? AND (index_next_at IS NULL OR index_next_at <= ?)", "pending", zeroDeletedAt, now).
		Order("id ASC").Limit(batch).Find(&docs).Error
	if err != nil {
		logx.Errorf("捞取待索引文档失败: %v", err)
		return
	}
	for i := range docs {
		claimAndProcess(ctx, svcCtx, &docs[i])
	}
}

// claimAndProcess 抢占式 UPDATE 把 pending 抢成 indexing：多副本只有一个抢到，天然去重。
func claimAndProcess(ctx context.Context, svcCtx *svc.ServiceContext, doc *knowledge_models.KnowledgeDocument) {
	res := svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ? AND index_status = ?", doc.DocumentID, "pending").
		Update("index_status", "indexing")
	if res.Error != nil || res.RowsAffected == 0 {
		return // 没抢到或出错，下一轮再战
	}
	if err := process(ctx, svcCtx, doc); err != nil {
		logx.Errorf("索引文档失败 document_id=%s: %v", doc.DocumentID, err)
		failDoc(svcCtx, doc, err)
		return
	}
	readyAt := time.Now()
	svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ?", doc.DocumentID).
		Updates(map[string]any{
			"index_status":      "ready",
			"index_retry_count": 0,
			"index_error":       "",
			"indexed_at":        &readyAt,
			"index_next_at":     nil,
		})
}

// process 分段→向量化→清旧索引→写新索引→回写 chunk 行。
func process(ctx context.Context, svcCtx *svc.ServiceContext, doc *knowledge_models.KnowledgeDocument) error {
	var coll knowledge_models.KnowledgeCollection
	if err := svcCtx.DB.Where("collection_id = ?", doc.CollectionID).First(&coll).Error; err != nil {
		return fmt.Errorf("加载集合失败 %s: %w", doc.CollectionID, err)
	}
	fallback := embedding.Config{
		ChunkSize: svcCtx.Config.Index.ChunkSize,
		Overlap:   svcCtx.Config.Index.Overlap,
	}
	ecfg := embedding.ParseIndexConfig(coll.IndexConfig, fallback)
	if !ecfg.Enabled() {
		return fmt.Errorf("集合 %s 未配置 embedding，无法索引", coll.CollectionID)
	}
	if ecfg.ChunkSize <= 0 {
		ecfg.ChunkSize = 800
	}

	// 正文未变且已 ready 的分段无需重算——但本函数只在 pending 时进入，一律重建。
	text := doc.ContentText
	pieces := segment([]rune(text), ecfg.ChunkSize, ecfg.Overlap)
	if len(pieces) == 0 {
		// 空正文：清掉旧分段即可，视为索引成功。
		return resetChunks(ctx, svcCtx, doc, nil, ecfg)
	}

	texts := make([]string, len(pieces))
	for i, p := range pieces {
		texts[i] = p
	}
	vectors, err := embedding.Embed(ctx, ecfg, texts)
	if err != nil {
		return fmt.Errorf("向量化失败: %w", err)
	}
	if len(vectors) != len(pieces) {
		return fmt.Errorf("向量数量不匹配: %d vs %d", len(vectors), len(pieces))
	}
	docs := make([]coreopensearch.ChunkDoc, 0, len(pieces))
	chunkRows := make([]knowledge_models.KnowledgeChunk, 0, len(pieces))
	for i, content := range pieces {
		docID := fmt.Sprintf("%s#%s#%d", doc.CollectionID, doc.DocumentID, i)
		docs = append(docs, coreopensearch.ChunkDoc{
			DocID:        docID,
			CollectionID: doc.CollectionID,
			DocumentID:   doc.DocumentID,
			Title:        doc.Title,
			ChunkIndex:   i,
			Content:      content,
			UserID:       doc.OwnerUserID,
			RagEnabled:   doc.RagEnabled,
			Embedding:    vectors[i],
		})
		chunkRows = append(chunkRows, knowledge_models.KnowledgeChunk{
			ChunkID:      "chk_" + shortHash(doc.DocumentID+fmt.Sprintf("-%d", i)),
			CollectionID: doc.CollectionID,
			DocumentID:   doc.DocumentID,
			ChunkIndex:   i,
			Content:      content,
			CharCount:    len([]rune(content)),
			EmbeddingRef: svcCtx.ChunkIndex + "#" + docID,
		})
	}
	return resetChunks(ctx, svcCtx, doc, docs, ecfg, chunkRows...)
}

// resetChunks 原子重建：先清该文档在 OpenSearch 与 DB 的旧分段，再写新分段。
// chunkRows 为空表示只清理（空正文）。
func resetChunks(ctx context.Context, svcCtx *svc.ServiceContext, doc *knowledge_models.KnowledgeDocument, osDocs []coreopensearch.ChunkDoc, ecfg embedding.Config, chunkRows ...knowledge_models.KnowledgeChunk) error {
	if svcCtx.OpenSearch != nil {
		if err := svcCtx.OpenSearch.DeleteByDocumentIDs(ctx, svcCtx.ChunkIndex, []string{doc.DocumentID}); err != nil {
			return fmt.Errorf("清理旧分段索引失败: %w", err)
		}
		if len(osDocs) > 0 {
			if err := svcCtx.OpenSearch.BulkUpsertChunks(ctx, svcCtx.ChunkIndex, osDocs); err != nil {
				return fmt.Errorf("写入分段索引失败: %w", err)
			}
		}
	}
	return svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("document_id = ?", doc.DocumentID).Delete(&knowledge_models.KnowledgeChunk{}).Error; err != nil {
			return err
		}
		if len(chunkRows) > 0 {
			if err := tx.CreateInBatches(chunkRows, 100).Error; err != nil {
				return err
			}
		}
		// 记住本正文指纹，供后续判重（本期不跨正文跳过，字段留作演进）。
		digest := sha256hex(doc.ContentText)
		return tx.Model(&knowledge_models.KnowledgeDocument{}).
			Where("document_id = ?", doc.DocumentID).
			Update("content_digest", digest).Error
	})
}

// failDoc 失败退避：未超上限则回 pending 并推迟 index_next_at，超限置 failed 并记 index_error。
func failDoc(svcCtx *svc.ServiceContext, doc *knowledge_models.KnowledgeDocument, cause error) {
	maxRetry := svcCtx.Config.Index.MaxRetry
	if maxRetry <= 0 {
		maxRetry = 3
	}
	next := doc.IndexRetryCount + 1
	msg := truncate(cause.Error(), 512)
	if next > maxRetry {
		svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
			Where("document_id = ?", doc.DocumentID).
			Updates(map[string]any{"index_status": "failed", "index_retry_count": next, "index_error": msg})
		return
	}
	// 指数退避：5s * 2^retry，封顶 30 分钟。
	backoff := time.Duration(5*(1<<uint(next))) * time.Second
	if backoff > 30*time.Minute {
		backoff = 30 * time.Minute
	}
	nextAt := time.Now().Add(backoff)
	svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ?", doc.DocumentID).
		Updates(map[string]any{"index_status": "pending", "index_retry_count": next, "index_error": msg, "index_next_at": &nextAt})
}

// segment 按字符窗口切段，窗口间保留 overlap 重叠，保住跨段语义连续。
func segment(r []rune, size, overlap int) []string {
	if size <= 0 {
		size = 800
	}
	if overlap < 0 || overlap >= size {
		overlap = 0
	}
	if len(r) == 0 {
		return nil
	}
	var out []string
	step := size - overlap
	for start := 0; start < len(r); start += step {
		end := start + size
		if end > len(r) {
			end = len(r)
		}
		piece := string(r[start:end])
		if len([]rune(piece)) > 0 {
			out = append(out, piece)
		}
		if end == len(r) {
			break
		}
	}
	return out
}

func sha256hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func shortHash(s string) string {
	return sha256hex(s)[:24]
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
