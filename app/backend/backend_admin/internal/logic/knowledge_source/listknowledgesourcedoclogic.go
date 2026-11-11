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
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListKnowledgeSourceDocLogic {
	return &ListKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("ListKnowledgeSourceDoc", ctx)}
}

// ListKnowledgeSourceDoc 列出知识源索引集合下的文档（在线文档引用与本地文件）。
// 知识源本身不存文档内容：文档本体与索引在 platform 的集合里，这里只做分页 / 关键字过滤。
func (l *ListKnowledgeSourceDocLogic) ListKnowledgeSourceDoc(req *types.ListKnowledgeSourceDocReq) (resp *types.ListKnowledgeSourceDocRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}
	row, err := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return &types.ListKnowledgeSourceDocRes{List: []types.KnowledgeSourceDoc{}, Total: 0}, nil
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
	list := make([]types.KnowledgeSourceDoc, 0, len(res.Documents))
	for _, doc := range res.Documents {
		// platform 的 ListDocuments 不支持关键字过滤，这里按标题做客户端过滤（文档量小，可接受）。
		if keyword != "" && !strings.Contains(strings.ToLower(doc.Title), strings.ToLower(keyword)) {
			continue
		}
		list = append(list, types.KnowledgeSourceDoc{
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
	return &types.ListKnowledgeSourceDocRes{List: list, Total: int64(len(list))}, nil
}
