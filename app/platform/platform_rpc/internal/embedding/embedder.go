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

package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config 集合 index_config 里的向量与分段配置。
// 单租户下 embedding 凭据直接配在知识库集合上（管理员录入），
// 与 platform AgentModel 的掩码接口解耦——worker 拿不到 AgentModel 的明文 key，
// 也就无需跨域读 platform 库：知识检索自包含在 knowledge 域内。
type Config struct {
	Endpoint  string `json:"embedding_endpoint"`
	APIKey    string `json:"embedding_api_key"`
	Model     string `json:"embedding_model"`
	Dims      int    `json:"embedding_dims"`
	ChunkSize int    `json:"chunk_size"`
	Overlap   int    `json:"overlap"`
	TopK      int    `json:"top_k"`
}

// Enabled 有 endpoint 且 model 才算可嵌入；缺凭据时检索/索引如实降级，不静默。
func (c Config) Enabled() bool {
	return strings.TrimSpace(c.Endpoint) != "" && strings.TrimSpace(c.Model) != ""
}

// ParseIndexConfig 解析集合 index_config JSON，缺项用 fallback 补齐。
// 坏 JSON 不致命：返回 fallback，让调用方按未配置处理并留日志。
func ParseIndexConfig(raw string, fallback Config) Config {
	cfg := fallback
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return cfg
	}
	var overlay Config
	if err := json.Unmarshal([]byte(raw), &overlay); err != nil {
		return cfg
	}
	if overlay.Endpoint != "" {
		cfg.Endpoint = overlay.Endpoint
	}
	if overlay.APIKey != "" {
		cfg.APIKey = overlay.APIKey
	}
	if overlay.Model != "" {
		cfg.Model = overlay.Model
	}
	if overlay.Dims > 0 {
		cfg.Dims = overlay.Dims
	}
	if overlay.ChunkSize > 0 {
		cfg.ChunkSize = overlay.ChunkSize
	}
	if overlay.Overlap > 0 {
		cfg.Overlap = overlay.Overlap
	}
	if overlay.TopK > 0 {
		cfg.TopK = overlay.TopK
	}
	return cfg
}

// Key 用于按 embedding 配置去重复用查询向量：同一 endpoint+model 才共享向量空间。
func (c Config) Key() string {
	return strings.TrimSpace(c.Endpoint) + "|" + strings.TrimSpace(c.Model)
}

// Embed 批量向量化。OpenAI 兼容 /embeddings 协议，一次请求一批文本，返回与入参等长的向量数组。
func Embed(ctx context.Context, cfg Config, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	if !cfg.Enabled() {
		return nil, fmt.Errorf("embedding 未配置（缺少 embedding_endpoint/embedding_model）")
	}
	payload, err := json.Marshal(map[string]any{"model": cfg.Model, "input": texts})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, embeddingsURL(cfg.Endpoint), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("embedding 上游返回 %d: %s", resp.StatusCode, truncate(string(body), 512))
	}
	var out struct {
		Data []struct {
			Index     int       `json:"index"`
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("embedding 响应解析失败: %w", err)
	}
	if len(out.Data) != len(texts) {
		return nil, fmt.Errorf("embedding 返回数量不匹配: 期望 %d 实得 %d", len(texts), len(out.Data))
	}
	// 按 index 落回原顺序，上游不保证返回有序。
	vectors := make([][]float32, len(texts))
	for _, d := range out.Data {
		if d.Index < 0 || d.Index >= len(texts) {
			return nil, fmt.Errorf("embedding 返回非法 index: %d", d.Index)
		}
		vectors[d.Index] = d.Embedding
	}
	return vectors, nil
}

// embeddingsURL 归一化到 /embeddings：兼容填 base（.../v1）与填完整 URL 两种写法。
func embeddingsURL(endpoint string) string {
	u := strings.TrimRight(strings.TrimSpace(endpoint), "/")
	if strings.HasSuffix(u, "/embeddings") {
		return u
	}
	return u + "/embeddings"
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
