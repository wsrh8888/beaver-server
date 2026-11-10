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

// ListDigitalKnowledgeSourceDoc 列出知识源索引集合下的文档（即用户粘贴的在线文档引用）。
// 知识源本身不存文档内容：文档本体与索引在 platform 的集合里，这里只做分页 / 关键字过滤。
func (l *ListDigitalKnowledgeSourceDocLogic) ListDigitalKnowledgeSourceDoc(req *types.ListDigitalKnowledgeSourceDocReq) (resp *types.ListDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return &types.ListDigitalKnowledgeSourceDocRes{List: []types.DigitalKnowledgeSourceDoc{}, Total: 0}, nil
	}

	res, rpcErr := l.svcCtx.PlatformRpc.ListDocuments(l.ctx, &platform.ListDocumentsReq{
		CollectionId: collectionID,
		Page:         1,
		PageSize:     200,
	})
	if rpcErr != nil {
		l.logger.Error(model.LogMsg{Text: "查询知识源文档失败", Data: map[string]interface{}{"err": rpcErr.Error(), "collectionId": collectionID}})
		return nil, errors.New("查询知识源文档失败")
	}
	keyword := strings.TrimSpace(req.Keyword)
	list := make([]types.DigitalKnowledgeSourceDoc, 0, len(res.Documents))
	for _, doc := range res.Documents {
		// platform 的 ListDocuments 不支持关键字过滤，这里按标题做客户端过滤（文档量小，可接受）。
		if keyword != "" && !strings.Contains(strings.ToLower(doc.Title), strings.ToLower(keyword)) {
			continue
		}
		list = append(list, types.DigitalKnowledgeSourceDoc{
			DocumentID:   doc.DocumentId,
			Title:        doc.Title,
			DocType:      doc.DocType,
			ByteSize:     doc.ByteSize,
			IndexStatus:  doc.IndexStatus,
			IndexError:   doc.IndexError,
			IndexedAt:    parseRFC3339Unix(doc.IndexedAt),
			SourceDocID:  doc.SourceDocId,
			SyncMode:     doc.SyncMode,
			LastSyncedAt: doc.LastSyncedAt,
		})
	}
	return &types.ListDigitalKnowledgeSourceDocRes{List: list, Total: int64(len(list))}, nil
}
