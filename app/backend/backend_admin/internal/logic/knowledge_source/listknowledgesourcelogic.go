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
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKnowledgeSourceLogic {
	return &ListKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("ListKnowledgeSource", ctx)}
}

// ListKnowledgeSource 列出某个场景下的知识源实体。
// 文档数是该知识源索引集合下的文档数（集合在新建知识源时创建，在线文档引用与本地文件都挂在上面）。
func (l *ListKnowledgeSourceLogic) ListKnowledgeSource(req *types.ListKnowledgeSourceReq) (resp *types.ListKnowledgeSourceRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}

	var rows []backend_models.AgentKnowledgeSource
	if err := l.svcCtx.DB.Where("scene = ?", scene).Order("id asc").Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询知识源失败", Data: map[string]interface{}{"err": err.Error(), "scene": scene}})
		return nil, errors.New("查询知识源失败")
	}

	list := make([]types.KnowledgeSourceItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.KnowledgeSourceItem{
			KnowledgeSourceID: row.SourceID,
			Name:              row.Name,
			RetrievalTopK:     row.RetrievalTopK,
			IndexConfig:       row.IndexConfig,
			DocumentCount:     collectionDocCount(l.ctx, l.svcCtx, l.logger, sourceCollectionID(row.KnowledgeIDs)),
			CreatedAt:         time.Time(row.CreatedAt).Unix(),
		})
	}
	return &types.ListKnowledgeSourceRes{List: list, Total: int64(len(list))}, nil
}
