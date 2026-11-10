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

type ListDigitalKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalKnowledgeSourceDocLogic {
	return &ListDigitalKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital_knowledge_source_doc", ctx)}
}

// ListDigitalKnowledgeSourceDoc 聚合展示知识源引入的所有知识库里的文档。
// 知识源本身不存文档：这里按引用逐个查 platform，再打上来源知识库标签。
// 单个知识库查询失败时跳过它，不因为一个库挂掉让整个详情打不开。
func (l *ListDigitalKnowledgeSourceDocLogic) ListDigitalKnowledgeSourceDoc(req *types.ListDigitalKnowledgeSourceDocReq) (resp *types.ListDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}

	names := collectionNameMap(l.ctx, l.svcCtx)
	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	list := make([]types.DigitalKnowledgeSourceDoc, 0)
	for _, collectionID := range uniqTrimmed(row.KnowledgeIDs) {
		res, rpcErr := l.svcCtx.PlatformRpc.ListDocuments(l.ctx, &platform.ListDocumentsReq{
			CollectionId: collectionID,
			Page:         1,
			PageSize:     200,
		})
		if rpcErr != nil {
			l.logger.Error(model.LogMsg{Text: "查询知识源文档失败", Data: map[string]interface{}{"err": rpcErr.Error(), "collectionId": collectionID}})
			continue
		}
		for _, doc := range res.Documents {
			if keyword != "" && !strings.Contains(strings.ToLower(doc.Title), keyword) {
				continue
			}
			list = append(list, types.DigitalKnowledgeSourceDoc{
				DocumentID:     doc.DocumentId,
				Title:          doc.Title,
				DocType:        doc.DocType,
				ByteSize:       doc.ByteSize,
				IndexStatus:    doc.IndexStatus,
				IndexError:     doc.IndexError,
				IndexedAt:      parseRFC3339Unix(doc.IndexedAt),
				CollectionID:   collectionID,
				CollectionName: names[collectionID],
			})
		}
	}
	return &types.ListDigitalKnowledgeSourceDocRes{List: list, Total: int64(len(list))}, nil
}
