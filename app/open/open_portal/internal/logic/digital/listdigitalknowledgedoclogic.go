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

type ListDigitalKnowledgeDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalKnowledgeDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalKnowledgeDocLogic {
	return &ListDigitalKnowledgeDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital_knowledge_doc", ctx)}
}

// ListDigitalKnowledgeDoc 知识库下的文档列表。
// platform 侧没有关键词检索参数，关键词过滤在这里做（文档量是单库级别，一次性拉取足够）。
func (l *ListDigitalKnowledgeDocLogic) ListDigitalKnowledgeDoc(req *types.ListDigitalKnowledgeDocReq) (resp *types.ListDigitalKnowledgeDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}
	if req.CollectionID == "" {
		return nil, errors.New("collectionId 不能为空")
	}

	res, err := l.svcCtx.PlatformRpc.ListDocuments(l.ctx, &platform.ListDocumentsReq{
		CollectionId: req.CollectionID,
		Page:         1,
		PageSize:     200,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询知识文档失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("查询文档失败")
	}

	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	list := make([]types.DigitalKnowledgeDoc, 0, len(res.Documents))
	for _, doc := range res.Documents {
		if keyword != "" && !strings.Contains(strings.ToLower(doc.Title), keyword) {
			continue
		}
		list = append(list, types.DigitalKnowledgeDoc{
			DocumentID:  doc.DocumentId,
			Title:       doc.Title,
			DocType:     doc.DocType,
			ByteSize:    doc.ByteSize,
			IndexStatus: doc.IndexStatus,
			IndexError:  doc.IndexError,
			IndexedAt:   parseRFC3339Unix(doc.IndexedAt),
		})
	}
	return &types.ListDigitalKnowledgeDocRes{List: list, Total: res.Total}, nil
}
