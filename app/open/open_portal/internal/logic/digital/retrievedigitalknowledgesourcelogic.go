/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"
	"strings"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type RetrieveDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRetrieveDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrieveDigitalKnowledgeSourceLogic {
	return &RetrieveDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("retrieve_digital_knowledge_source", ctx)}
}

// RetrieveDigitalKnowledgeSource 对该知识源索引集合做一次召回测试。
// 用户添加的在线文档引用都挂在这个集合上，检索口径与运行时一致。
func (l *RetrieveDigitalKnowledgeSourceLogic) RetrieveDigitalKnowledgeSource(req *types.RetrieveDigitalKnowledgeSourceReq) (resp *types.RetrieveDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return &types.RetrieveDigitalKnowledgeSourceRes{Hits: []types.DigitalKnowledgeSourceHit{}, Notice: "请先为这个知识源添加在线文档"}, nil
	}
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return &types.RetrieveDigitalKnowledgeSourceRes{Hits: []types.DigitalKnowledgeSourceHit{}, Notice: "请输入测试问题"}, nil
	}
	topK := req.TopK
	if topK <= 0 {
		topK = row.RetrievalTopK
	}
	if topK <= 0 || topK > open_models.MaxKnowledgeSourceTopK {
		topK = open_models.DefaultKnowledgeSourceTopK
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
	hits := make([]types.DigitalKnowledgeSourceHit, 0, len(res.Hits))
	for _, hit := range res.Hits {
		hits = append(hits, types.DigitalKnowledgeSourceHit{
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
	return &types.RetrieveDigitalKnowledgeSourceRes{Hits: hits, Notice: notice}, nil
}
