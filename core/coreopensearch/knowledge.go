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

package coreopensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 知识检索相关能力：索引 mapping、批量写入、按文档删除、knn 向量检索。
// 与 opensearch.go 的日志写入共用同一个极简 HTTP Client，但独立成文件，
// 避免把 RAG 的复杂查询体混进「日志扁平写入」这条稳定路径。

// EmbeddingField 向量字段名，mapping 与 knn 查询必须一致。
const EmbeddingField = "embedding"

// ChunkDoc OpenSearch 里一条分段文档。向量本体存这里，DB 只存正文+段号。
type ChunkDoc struct {
	DocID        string    `json:"-"` // _id：collection#document#chunk_index，幂等覆盖
	CollectionID string    `json:"collection_id"`
	DocumentID   string    `json:"document_id"`
	Title        string    `json:"title"`
	ChunkIndex   int       `json:"chunk_index"`
	Content      string    `json:"content"`
	UserID       string    `json:"user_id,omitempty"`
	RagEnabled   bool      `json:"rag_enabled"`
	Embedding    []float32 `json:"embedding"`
}

// SearchHit 一条 knn 命中。
type SearchHit struct {
	CollectionID string
	DocumentID   string
	Title        string
	ChunkIndex   int
	Content      string
	Score        float64
}

// EnsureKnowledgeIndex 创建知识分段索引（幂等：已存在直接返回）。
// dims 为 embedding 模型输出维度，knn_vector 建索引时维度定死，改模型须重建索引。
func (c *Client) EnsureKnowledgeIndex(ctx context.Context, index string, dims int) error {
	if c == nil || c.http == nil {
		return fmt.Errorf("OpenSearch 未初始化")
	}
	index = strings.TrimSpace(index)
	if index == "" {
		return fmt.Errorf("OpenSearch index 为空")
	}
	body := map[string]any{
		"settings": map[string]any{"index.knn": true},
		"mappings": map[string]any{
			"properties": map[string]any{
				EmbeddingField: map[string]any{
					"type":      "knn_vector",
					"dimension": dims,
					"method": map[string]any{
						"name":       "hnsw",
						"space_type": "cosinesimil",
						"engine":     "nmslib",
						"parameters": map[string]any{"ef_construction": 128, "m": 16},
					},
				},
				"collection_id": map[string]any{"type": "keyword"},
				"document_id":   map[string]any{"type": "keyword"},
				"title":         map[string]any{"type": "keyword"},
				"chunk_index":   map[string]any{"type": "integer"},
				"content":       map[string]any{"type": "text"},
				"user_id":       map[string]any{"type": "keyword"},
				"rag_enabled":   map[string]any{"type": "boolean"},
			},
		},
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/"+url.PathEscape(index), bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	// 400 resource_already_exists_exception 视为成功：幂等建库。
	if resp.StatusCode >= 300 && !strings.Contains(string(respBody), "resource_already_exists_exception") {
		return fmt.Errorf("opensearch create index failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}

// BulkUpsertChunks 批量写入分段（_id 幂等覆盖，重复索引不产生脏数据）。
func (c *Client) BulkUpsertChunks(ctx context.Context, index string, docs []ChunkDoc) error {
	if c == nil || c.http == nil {
		return fmt.Errorf("OpenSearch 未初始化")
	}
	if len(docs) == 0 {
		return nil
	}
	index = strings.TrimSpace(index)
	var ndjson strings.Builder
	for _, d := range docs {
		meta := map[string]any{"index": map[string]any{"_index": index, "_id": d.DocID}}
		mb, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		db, err := json.Marshal(d)
		if err != nil {
			return err
		}
		ndjson.Write(mb)
		ndjson.WriteByte('\n')
		ndjson.Write(db)
		ndjson.WriteByte('\n')
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/_bulk", strings.NewReader(ndjson.String()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("opensearch bulk failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	// bulk 是部分成功语义：HTTP 200 仍可能有个别 item 失败。
	var out struct {
		Errors bool `json:"errors"`
	}
	_ = json.Unmarshal(respBody, &out)
	if out.Errors {
		return fmt.Errorf("opensearch bulk 部分写入失败: %s", string(respBody))
	}
	return nil
}

// DeleteByDocumentIDs 按文档 ID 批量清除其全部分段（删文档/停用/重建前置清理）。
func (c *Client) DeleteByDocumentIDs(ctx context.Context, index string, documentIDs []string) error {
	if c == nil || c.http == nil {
		return fmt.Errorf("OpenSearch 未初始化")
	}
	if len(documentIDs) == 0 {
		return nil
	}
	index = strings.TrimSpace(index)
	body, err := json.Marshal(map[string]any{
		"query": map[string]any{"terms": map[string]any{"document_id": documentIDs}},
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+url.PathEscape(index)+"/_delete_by_query", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	// 索引不存在时删除无意义：404 视为成功（本就没有分段）。
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("opensearch delete_by_query failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}

// KnnSearch 在指定集合范围内做向量检索。collectionIDs 为空直接返回空（防越权：不给范围就不查）。
// filter 强制 rag_enabled=true，停用/不参与检索的分段永不命中。
func (c *Client) KnnSearch(ctx context.Context, index string, vector []float32, collectionIDs []string, k int) ([]SearchHit, error) {
	if c == nil || c.http == nil {
		return nil, fmt.Errorf("OpenSearch 未初始化")
	}
	index = strings.TrimSpace(index)
	if len(vector) == 0 || len(collectionIDs) == 0 || k <= 0 {
		return nil, nil
	}
	query := map[string]any{
		"size": k,
		"query": map[string]any{
			"knn": map[string]any{
				EmbeddingField: map[string]any{
					"vector": vector,
					"k":      k,
					"filter": map[string]any{
						"bool": map[string]any{
							"must": []any{
								map[string]any{"terms": map[string]any{"collection_id": collectionIDs}},
								map[string]any{"term": map[string]any{"rag_enabled": true}},
							},
						},
					},
				},
			},
		},
	}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+url.PathEscape(index)+"/_search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("opensearch knn failed: status=%d body=%s", resp.StatusCode, string(raw))
	}
	var parsed struct {
		Hits struct {
			Hits []struct {
				Score  float64 `json:"_score"`
				Source struct {
					CollectionID string `json:"collection_id"`
					DocumentID   string `json:"document_id"`
					Title        string `json:"title"`
					ChunkIndex   int    `json:"chunk_index"`
					Content      string `json:"content"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	hits := make([]SearchHit, 0, len(parsed.Hits.Hits))
	for _, h := range parsed.Hits.Hits {
		hits = append(hits, SearchHit{
			CollectionID: h.Source.CollectionID,
			DocumentID:   h.Source.DocumentID,
			Title:        h.Source.Title,
			ChunkIndex:   h.Source.ChunkIndex,
			Content:      h.Source.Content,
			Score:        h.Score,
		})
	}
	return hits, nil
}
