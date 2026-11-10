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

	"beaver/app/doc/doc_rpc/doc"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ImportDigitalKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewImportDigitalKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportDigitalKnowledgeSourceDocLogic {
	return &ImportDigitalKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("import_digital_knowledge_source_doc", ctx)}
}

// ImportDigitalKnowledgeSourceDoc 把用户粘贴的海狸内部在线文档地址导入到知识源。
//
// 这里不做任何正文搬运：先按 user_id 调 doc_rpc.GetDoc 鉴权（文档不可见 / 无权限直接计入失败），
// 再用 source_type=doc 调 platform_rpc.CreateDocument 建立引用，正文与索引由 platform 的
// ingest 异步从在线文档拉取（飞书/钉钉「引用云文档」语义）。在线文档本体永远留在 doc 域，
// 知识源这一侧只存引用与索引产物。
func (l *ImportDigitalKnowledgeSourceDocLogic) ImportDigitalKnowledgeSourceDoc(req *types.ImportDigitalKnowledgeSourceDocReq) (resp *types.ImportDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID)
	if err != nil {
		return nil, err
	}
	collectionID := sourceCollectionID(row.KnowledgeIDs)
	if collectionID == "" {
		return nil, errors.New("知识源尚未就绪，请先保存知识源")
	}

	urls := req.DocURLs
	if len(urls) == 0 {
		return &types.ImportDigitalKnowledgeSourceDocRes{Imported: 0, Failed: []string{}, Notice: "没有需要导入的地址"}, nil
	}

	imported := 0
	failed := make([]string, 0, len(urls))
	for _, raw := range urls {
		docID := parseBeaverDocID(raw)
		if docID == "" {
			failed = append(failed, raw)
			continue
		}
		// 1) 按当前用户鉴权：文档不可见 / 无权限则 GetDoc 报错，计入失败。
		if _, authErr := l.svcCtx.DocRpc.GetDoc(l.ctx, &doc.GetDocReq{
			DocId:  docID,
			UserId: req.UserID,
		}); authErr != nil {
			l.logger.Warn(model.LogMsg{Text: "在线文档鉴权失败，跳过导入", Data: map[string]interface{}{"docId": docID, "err": authErr.Error()}})
			failed = append(failed, raw)
			continue
		}
		// 2) 建立引用：platform 侧 ingest 会按 doc_id 自动拉取并索引源文档。
		//    direct 模式表示跟随源文档更新（引用云文档），正文与索引由平台侧异步完成。
		if _, createErr := l.svcCtx.PlatformRpc.CreateDocument(l.ctx, &platform.CreateDocumentReq{
			CollectionId: collectionID,
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
		notice = "全部地址解析或鉴权失败，未导入任何文档"
	} else if len(failed) > 0 {
		notice = "部分地址解析或鉴权失败，已跳过"
	}
	return &types.ImportDigitalKnowledgeSourceDocRes{
		Imported: imported,
		Failed:   failed,
		Notice:   notice,
	}, nil
}
