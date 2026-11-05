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

package logic

import (
	"context"
	"sort"

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/embedding"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"
	"beaver/core/coreopensearch"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchKnowledgeLogic {
	return &SearchKnowledgeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// SearchKnowledge 运行时向量检索：只认 app 绑定的知识库，不接受调用方指定分区。
// 无绑定 / 未命中 / OpenSearch 未配置 / embedding 未配，一律返回空命中不报错——
// 对齐 stratum「没查到就用空结果，别拿常识补全」的降级语义。
func (l *SearchKnowledgeLogic) SearchKnowledge(in *knowledge_rpc.SearchKnowledgeReq) (*knowledge_rpc.SearchKnowledgeRes, error) {
	if trim(in.Query) == "" {
		return &knowledge_rpc.SearchKnowledgeRes{}, nil
	}
	bound := resolveBoundCollectionIDs(l.ctx, l.svcCtx, in.AppId, in.Scene)
	if len(bound) == 0 {
		return &knowledge_rpc.SearchKnowledgeRes{}, nil
	}
	if l.svcCtx.OpenSearch == nil {
		l.Errorf("知识检索降级：OpenSearch 未配置")
		return &knowledge_rpc.SearchKnowledgeRes{}, nil
	}

	// 只在「启用中的绑定集合」里检索：停用集合的旧分段即便留在索引里也不再命中。
	var colls []knowledge_models.KnowledgeCollection
	if err := l.svcCtx.DB.Where("collection_id IN ? AND status = ?", bound, 1).Find(&colls).Error; err != nil {
		l.Errorf("加载绑定集合失败: %v", err)
		return nil, err
	}
	if len(colls) == 0 {
		return &knowledge_rpc.SearchKnowledgeRes{}, nil
	}

	topK := int(in.TopK)
	if topK <= 0 {
		topK = 5
	}

	// 按 embedding 配置分组：同一 endpoint+model 共享向量空间，查询向量只算一次。
	fallback := embedding.Config{Dims: l.svcCtx.Config.Index.Dims}
	type embGroup struct {
		cfg embedding.Config
		ids []string
	}
	groups := map[string]*embGroup{}
	order := []string{}
	for i := range colls {
		cfg := embedding.ParseIndexConfig(colls[i].IndexConfig, fallback)
		if !cfg.Enabled() {
			continue // 未配 embedding 的集合无法做向量检索，如实跳过
		}
		key := cfg.Key()
		g, ok := groups[key]
		if !ok {
			g = &embGroup{cfg: cfg}
			groups[key] = g
			order = append(order, key)
		}
		g.ids = append(g.ids, colls[i].CollectionID)
	}

	var all []coreopensearch.SearchHit
	for _, key := range order {
		g := groups[key]
		vecs, err := embedding.Embed(l.ctx, g.cfg, []string{in.Query})
		if err != nil {
			l.Errorf("查询向量化失败(%s): %v", key, err)
			continue
		}
		if len(vecs) == 0 || len(vecs[0]) == 0 {
			continue
		}
		hits, err := l.svcCtx.OpenSearch.KnnSearch(l.ctx, l.svcCtx.ChunkIndex, vecs[0], g.ids, topK)
		if err != nil {
			l.Errorf("knn 检索失败(%s): %v", key, err)
			continue
		}
		all = append(all, hits...)
	}

	// 跨分组按相关度合并，裁剪到 top_k。
	sort.SliceStable(all, func(i, j int) bool { return all[i].Score > all[j].Score })
	if len(all) > topK {
		all = all[:topK]
	}
	hits := make([]*knowledge_rpc.KnowledgeHit, 0, len(all))
	for _, h := range all {
		hits = append(hits, &knowledge_rpc.KnowledgeHit{
			DocumentId:   h.DocumentID,
			Title:        h.Title,
			ChunkIndex:   int32(h.ChunkIndex),
			Content:      h.Content,
			Score:        h.Score,
			CollectionId: h.CollectionID,
		})
	}
	return &knowledge_rpc.SearchKnowledgeRes{Hits: hits}, nil
}
