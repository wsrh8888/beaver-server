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

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/embedding"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetrieveCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRetrieveCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrieveCollectionLogic {
	return &RetrieveCollectionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// RetrieveCollection 管理端召回测试：按集合直接检索。
// 与 SearchKnowledge 的区别只在「谁能查哪个集合」——运行时只认 app 绑定，
// 这里是控制台调试入口，允许指定单个集合；除此以外降级语义完全一致：
// 未配 OpenSearch / 集合未配向量检索 / 无命中，一律返回空命中 + notice 说明原因，不报错。
func (l *RetrieveCollectionLogic) RetrieveCollection(in *platform_rpc.RetrieveCollectionReq) (*platform_rpc.RetrieveCollectionRes, error) {
	query := trim(in.Query)
	if query == "" {
		return &platform_rpc.RetrieveCollectionRes{Notice: "请输入测试问题"}, nil
	}
	collectionID := trim(in.CollectionId)
	if collectionID == "" {
		return &platform_rpc.RetrieveCollectionRes{Notice: "未指定知识库"}, nil
	}

	var coll knowledge_models.KnowledgeCollection
	if err := l.svcCtx.DB.Where("collection_id = ?", collectionID).First(&coll).Error; err != nil {
		return &platform_rpc.RetrieveCollectionRes{Notice: "知识库不存在"}, nil
	}
	if coll.Status != 1 {
		return &platform_rpc.RetrieveCollectionRes{Notice: "知识库已停用，检索不会命中"}, nil
	}

	if l.svcCtx.OpenSearch == nil {
		l.Errorf("召回测试降级：OpenSearch 未配置")
		return &platform_rpc.RetrieveCollectionRes{Notice: "服务端未配置 OpenSearch，无法检索"}, nil
	}

	cfg := embedding.ParseIndexConfig(coll.IndexConfig, embedding.Config{Dims: l.svcCtx.Config.Index.Dims})
	if !cfg.Enabled() {
		return &platform_rpc.RetrieveCollectionRes{Notice: "该知识库未配置向量模型，无法做召回测试"}, nil
	}

	topK := int(in.TopK)
	if topK <= 0 {
		topK = 5
	}

	vecs, err := embedding.Embed(l.ctx, cfg, []string{query})
	if err != nil {
		l.Errorf("查询向量化失败: %v", err)
		return &platform_rpc.RetrieveCollectionRes{Notice: "查询向量化失败，请检查向量模型配置"}, nil
	}
	if len(vecs) == 0 || len(vecs[0]) == 0 {
		return &platform_rpc.RetrieveCollectionRes{Notice: "查询向量化返回空，请检查向量模型配置"}, nil
	}

	hits, err := l.svcCtx.OpenSearch.KnnSearch(l.ctx, l.svcCtx.ChunkIndex, vecs[0], []string{coll.CollectionID}, topK)
	if err != nil {
		l.Errorf("knn 检索失败: %v", err)
		return &platform_rpc.RetrieveCollectionRes{Notice: "检索失败，请查看服务端日志"}, nil
	}
	if len(hits) == 0 {
		return &platform_rpc.RetrieveCollectionRes{Notice: "无命中：文档可能还没学完，或与问题不相关"}, nil
	}

	sort.SliceStable(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	if len(hits) > topK {
		hits = hits[:topK]
	}
	list := make([]*platform_rpc.KnowledgeHit, 0, len(hits))
	for _, h := range hits {
		list = append(list, &platform_rpc.KnowledgeHit{
			DocumentId:   h.DocumentID,
			Title:        h.Title,
			ChunkIndex:   int32(h.ChunkIndex),
			Content:      h.Content,
			Score:        h.Score,
			CollectionId: h.CollectionID,
		})
	}
	return &platform_rpc.RetrieveCollectionRes{Hits: list}, nil
}
