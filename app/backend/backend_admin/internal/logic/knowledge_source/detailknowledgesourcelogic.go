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
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type DetailKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDetailKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailKnowledgeSourceLogic {
	return &DetailKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("DetailKnowledgeSource", ctx)}
}

// DetailKnowledgeSource 按 id 取回单个知识源，供详情子路由直接进页时加载。
func (l *DetailKnowledgeSourceLogic) DetailKnowledgeSource(req *types.DetailKnowledgeSourceReq) (resp *types.DetailKnowledgeSourceRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}
	row, err := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	return &types.DetailKnowledgeSourceRes{
		Source: types.KnowledgeSourceItem{
			KnowledgeSourceID: row.SourceID,
			Name:              row.Name,
			RetrievalTopK:     row.RetrievalTopK,
			IndexConfig:       row.IndexConfig,
			DocumentCount:     collectionDocCount(l.ctx, l.svcCtx, l.logger, sourceCollectionID(row.KnowledgeIDs)),
			CreatedAt:         time.Time(row.CreatedAt).Unix(),
		},
	}, nil
}
