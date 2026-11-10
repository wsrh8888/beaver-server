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
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListDigitalKnowledgeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalKnowledgeLogic {
	return &ListDigitalKnowledgeLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital_knowledge", ctx)}
}

// ListDigitalKnowledge 列出平台知识库；传入 digitalId 时额外标出哪些已绑定该数字伙伴。
// 文档数由 platform 逐集合统计：集合表本身不存计数，跨库也读不到 knowledge_documents。
func (l *ListDigitalKnowledgeLogic) ListDigitalKnowledge(req *types.ListDigitalKnowledgeReq) (resp *types.ListDigitalKnowledgeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.DigitalID != "" {
		if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
			return nil, err
		}
	}

	res, err := l.svcCtx.PlatformRpc.ListCollections(l.ctx, &platform.ListCollectionsReq{Page: 1, PageSize: 200})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询知识库列表失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("查询知识库失败")
	}

	bound := map[string]bool{}
	if req.DigitalID != "" {
		for _, id := range loadDigitalKnowledgeIDs(l.svcCtx.DB, req.DigitalID) {
			bound[id] = true
		}
	}

	list := make([]types.DigitalKnowledgeItem, 0, len(res.Collections))
	for _, item := range res.Collections {
		docCount := int64(0)
		if doc, err := l.svcCtx.PlatformRpc.ListDocuments(l.ctx, &platform.ListDocumentsReq{
			CollectionId: item.CollectionId,
			Page:         1,
			PageSize:     1,
		}); err == nil {
			docCount = doc.Total
		}
		list = append(list, types.DigitalKnowledgeItem{
			CollectionID:  item.CollectionId,
			Name:          item.Name,
			Description:   item.Description,
			Kind:          item.Kind,
			Status:        int(item.Status),
			TopK:          indexConfigTopK(item.IndexConfig),
			DocumentCount: docCount,
			Bound:         bound[item.CollectionId],
			IndexConfig:   item.IndexConfig,
		})
	}
	return &types.ListDigitalKnowledgeRes{List: list, Total: res.Total}, nil
}
