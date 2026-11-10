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

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListDigitalKnowledgeSourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalKnowledgeSourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalKnowledgeSourceLogic {
	return &ListDigitalKnowledgeSourceLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital_knowledge_source", ctx)}
}

// ListDigitalKnowledgeSource 列出某个数字伙伴下的知识源实体。
// 文档数是它引入的所有知识库的文档数之和：知识源本身不存文档，只是引用分组。
func (l *ListDigitalKnowledgeSourceLogic) ListDigitalKnowledgeSource(req *types.ListDigitalKnowledgeSourceReq) (resp *types.ListDigitalKnowledgeSourceRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}
	if err := migrateLegacyKnowledgeSources(l.svcCtx.DB, req.DigitalID, req.UserID); err != nil {
		l.logger.Error(model.LogMsg{Text: "迁移旧知识源配置失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("读取知识源失败")
	}

	var rows []open_models.OpenKnowledgeSource
	if err := l.svcCtx.DB.Where("digital_id = ?", req.DigitalID).Order("id asc").Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询知识源失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("查询知识源失败")
	}

	list := make([]types.DigitalKnowledgeSourceItem, 0, len(rows))
	for _, row := range rows {
		ids := append([]string{}, row.KnowledgeIDs...)
		list = append(list, types.DigitalKnowledgeSourceItem{
			KnowledgeSourceID: row.SourceID,
			Name:              row.Name,
			RetrievalTopK:     row.RetrievalTopK,
			IndexConfig:       row.IndexConfig,
			KnowledgeIDs:      ids,
			DocumentCount:     countCollectionDocuments(l.ctx, l.svcCtx, ids),
			CreatedAt:         row.CreatedAt.Unix(),
		})
	}
	return &types.ListDigitalKnowledgeSourceRes{List: list, Total: int64(len(list))}, nil
}
