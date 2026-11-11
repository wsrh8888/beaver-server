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
	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ImportKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewImportKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportKnowledgeSourceDocLogic {
	return &ImportKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("ImportKnowledgeSourceDoc", ctx)}
}

// ImportKnowledgeSourceDoc 把运营粘贴的海狸内部在线文档地址导入到知识源。
//
// 这里不做任何正文搬运：先按 doc_id 读取快照确认文档可用，再用 source_type=doc 建引用，
// 正文与索引由 platform 的 ingest 从在线文档拉取（引用云文档语义）。在线文档本体永远留在 doc 域。
func (l *ImportKnowledgeSourceDocLogic) ImportKnowledgeSourceDoc(req *types.ImportKnowledgeSourceDocReq) (resp *types.ImportKnowledgeSourceDocRes, err error) {
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
		return nil, errors.New("知识源尚未就绪，请先保存知识源")
	}

	urls := req.DocURLs
	if len(urls) == 0 {
		return &types.ImportKnowledgeSourceDocRes{Imported: 0, Failed: []string{}, Notice: "没有需要导入的地址"}, nil
	}

	docType := strings.ToLower(strings.TrimSpace(req.DocType))
	if docType == "" {
		docType = "doc"
	}

	imported := 0
	failed := make([]string, 0, len(urls))
	for _, raw := range urls {
		docID := parseBeaverDocID(raw)
		if docID == "" {
			failed = append(failed, raw)
			continue
		}
		// 1) 确认文档可读：读不到快照就不建引用，避免留下永远索引不出来的空引用。
		if _, readErr := l.svcCtx.DocRpc.GetDocContent(l.ctx, &doccli.GetDocContentReq{DocId: docID}); readErr != nil {
			l.logger.Warn(model.LogMsg{Text: "在线文档读取失败，跳过导入", Data: map[string]interface{}{"docId": docID, "err": readErr.Error()}})
			failed = append(failed, raw)
			continue
		}
		// 2) 建立引用：platform 侧 ingest 会按 doc_id 自动拉取并索引源文档。
		if _, createErr := l.svcCtx.PlatformRpc.CreateDocument(l.ctx, &platform.CreateDocumentReq{
			CollectionId: collectionID,
			DocType:      docType,
			SourceType:   "doc",
			SourceDocId:  docID,
			SyncMode:     "direct",
			RagEnabled:   true,
		}); createErr != nil {
			l.logger.Error(model.LogMsg{Text: "建立在线文档引用失败", Data: map[string]interface{}{"docId": docID, "err": createErr.Error()}})
			failed = append(failed, raw)
			continue
		}
		imported++
	}

	notice := ""
	if imported == 0 {
		notice = "全部地址解析或读取失败，未导入任何文档"
	} else if len(failed) > 0 {
		notice = "部分地址解析或读取失败，已跳过"
	}
	return &types.ImportKnowledgeSourceDocRes{Imported: imported, Failed: failed, Notice: notice}, nil
}
