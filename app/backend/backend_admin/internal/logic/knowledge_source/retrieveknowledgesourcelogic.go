/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package knowledge_source

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type RetrieveKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRetrieveKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrieveKnowledgeSourceLogic {
	return &RetrieveKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("RetrieveKnowledgeSource", ctx)}
}

// RetrieveKnowledgeSource 对该知识源索引集合做一次召回测试。
// 文档引用都挂在这个集合上，检索口径与运行时一致。
func (l *RetrieveKnowledgeSourceLogic) RetrieveKnowledgeSource(req *types.RetrieveKnowledgeSourceReq) (resp *types.RetrieveKnowledgeSourceRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}
	row, err := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return &types.RetrieveKnowledgeSourceRes{Hits: []types.KnowledgeSourceHit{}, Notice: "请先为这个知识源添加文档"}, nil
	}
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return &types.RetrieveKnowledgeSourceRes{Hits: []types.KnowledgeSourceHit{}, Notice: "请输入测试问题"}, nil
	}
	topK := req.TopK
	if topK <= 0 {
		topK = row.RetrievalTopK
	}
	if topK <= 0 || topK > backend_models.MaxAgentKnowledgeSourceTopK {
		topK = backend_models.DefaultAgentKnowledgeSourceTopK
	}

	res, rpcErr := l.svcCtx.PlatformRpc.RetrieveCollection(l.ctx, &platform.RetrieveCollectionReq{
		CollectionId: collectionID,
		Query:        query,
		TopK:         int32(topK),
	})
	if rpcErr != nil {
		l.logger.Error(model.LogMsg{Text: "知识源召回测试失败", Data: map[string]interface{}{"err": rpcErr.Error(), "collectionId": collectionID}})
		return nil, errors.New("召回测试失败")
	}
	hits := make([]types.KnowledgeSourceHit, 0, len(res.Hits))
	for _, hit := range res.Hits {
		hits = append(hits, types.KnowledgeSourceHit{
			DocumentID: hit.DocumentId,
			Title:      hit.Title,
			ChunkIndex: int(hit.ChunkIndex),
			Content:    hit.Content,
			Score:      hit.Score,
		})
	}
	notice := res.Notice
	if len(hits) == 0 && notice == "" {
		notice = "未召回相关内容"
	}
	return &types.RetrieveKnowledgeSourceRes{Hits: hits, Notice: notice}, nil
}
