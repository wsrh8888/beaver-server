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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type RetrieveDigitalKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRetrieveDigitalKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetrieveDigitalKnowledgeLogic {
	return &RetrieveDigitalKnowledgeLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("retrieve_digital_knowledge", ctx)}
}

// RetrieveDigitalKnowledge 召回测试：只查本数字伙伴已绑定的知识库，按集合单独检索。
// 未命中 / 未配向量模型都返回 notice 说明，不报错——控制台要能解释「为什么没结果」。
func (l *RetrieveDigitalKnowledgeLogic) RetrieveDigitalKnowledge(req *types.RetrieveDigitalKnowledgeReq) (resp *types.RetrieveDigitalKnowledgeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return &types.RetrieveDigitalKnowledgeRes{Hits: []types.DigitalKnowledgeHit{}, Notice: "请输入测试问题"}, nil
	}

	bound := false
	for _, id := range loadDigitalKnowledgeIDs(l.svcCtx.DB, partner.DigitalID) {
		if id == req.CollectionID {
			bound = true
			break
		}
	}
	if !bound {
		return &types.RetrieveDigitalKnowledgeRes{Hits: []types.DigitalKnowledgeHit{}, Notice: "该知识库还没有绑定到当前数字伙伴"}, nil
	}

	res, err := l.svcCtx.PlatformRpc.RetrieveCollection(l.ctx, &platform.RetrieveCollectionReq{
		CollectionId: req.CollectionID,
		Query:        query,
		TopK:         int32(req.TopK),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "召回测试失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("召回失败")
	}

	hits := make([]types.DigitalKnowledgeHit, 0, len(res.Hits))
	for _, hit := range res.Hits {
		hits = append(hits, types.DigitalKnowledgeHit{
			DocumentID: hit.DocumentId,
			Title:      hit.Title,
			ChunkIndex: int(hit.ChunkIndex),
			Content:    hit.Content,
			Score:      hit.Score,
		})
	}
	return &types.RetrieveDigitalKnowledgeRes{Hits: hits, Notice: res.Notice}, nil
}
