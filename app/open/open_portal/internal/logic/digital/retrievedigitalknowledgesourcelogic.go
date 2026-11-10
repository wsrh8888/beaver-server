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
	"sort"
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

// RetrieveDigitalKnowledgeSource 跨知识源引入的所有知识库做召回测试。
// 每个库各自检索 topK，再按相关度统一排序截断，保证「一个知识源引入 N 个库」
// 的检索结果与运行时按扁平 knowledgeIds 检索的口径一致。
func (l *RetrieveDigitalKnowledgeSourceLogic) RetrieveDigitalKnowledgeSource(req *types.RetrieveDigitalKnowledgeSourceReq) (resp *types.RetrieveDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return &types.RetrieveDigitalKnowledgeSourceRes{Hits: []types.DigitalKnowledgeSourceHit{}, Notice: "请输入测试问题"}, nil
	}
	ids := uniqTrimmed(row.KnowledgeIDs)
	if len(ids) == 0 {
		return &types.RetrieveDigitalKnowledgeSourceRes{Hits: []types.DigitalKnowledgeSourceHit{}, Notice: "请先为这个知识源引入知识库"}, nil
	}

	topK := req.TopK
	if topK <= 0 {
		topK = row.RetrievalTopK
	}
	if topK <= 0 || topK > open_models.MaxKnowledgeSourceTopK {
		topK = open_models.DefaultKnowledgeSourceTopK
	}

	names := collectionNameMap(l.ctx, l.svcCtx)
	hits := make([]types.DigitalKnowledgeSourceHit, 0)
	notice := ""
	for _, collectionID := range ids {
		res, rpcErr := l.svcCtx.PlatformRpc.RetrieveCollection(l.ctx, &platform.RetrieveCollectionReq{
			CollectionId: collectionID,
			Query:        query,
			TopK:         int32(topK),
		})
		if rpcErr != nil {
			l.logger.Error(model.LogMsg{Text: "知识源召回测试失败", Data: map[string]interface{}{"err": rpcErr.Error(), "collectionId": collectionID}})
			continue
		}
		if notice == "" {
			notice = res.Notice
		}
		for _, hit := range res.Hits {
			hits = append(hits, types.DigitalKnowledgeSourceHit{
				DocumentID:     hit.DocumentId,
				CollectionID:   collectionID,
				CollectionName: names[collectionID],
				Title:          hit.Title,
				ChunkIndex:     int(hit.ChunkIndex),
				Content:        hit.Content,
				Score:          hit.Score,
			})
		}
	}
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	if len(hits) > topK {
		hits = hits[:topK]
	}
	if len(hits) == 0 && notice == "" {
		notice = "未召回相关内容"
	}
	return &types.RetrieveDigitalKnowledgeSourceRes{Hits: hits, Notice: notice}, nil
}
