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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type DetailDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDetailDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailDigitalKnowledgeSourceLogic {
	return &DetailDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("detail_digital_knowledge_source", ctx)}
}

// DetailDigitalKnowledgeSource 按 id 取回单个知识源实体，供子路由进入时加载。
// 归属校验走 loadOwnedKnowledgeSource：跨数字伙伴不能读别人的知识源。
func (l *DetailDigitalKnowledgeSourceLogic) DetailDigitalKnowledgeSource(req *types.DetailDigitalKnowledgeSourceReq) (resp *types.DetailDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	source := types.DigitalKnowledgeSourceItem{
		KnowledgeSourceID: row.SourceID,
		Name:              row.Name,
		RetrievalTopK:     row.RetrievalTopK,
		IndexConfig:       row.IndexConfig,
		DocumentCount:     collectionDocCount(l.ctx, l.svcCtx, l.logger, sourceCollectionID(row.KnowledgeIDs)),
		CreatedAt:         row.CreatedAt.Unix(),
	}
	return &types.DetailDigitalKnowledgeSourceRes{Source: source}, nil
}
