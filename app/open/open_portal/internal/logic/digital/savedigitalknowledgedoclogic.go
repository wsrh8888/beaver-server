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

type SaveDigitalKnowledgeDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDigitalKnowledgeDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalKnowledgeDocLogic {
	return &SaveDigitalKnowledgeDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_digital_knowledge_doc", ctx)}
}

// SaveDigitalKnowledgeDoc 新建或更新知识文档。
// 正文变更由 platform 侧判重后自动重新入队索引，此处不需要额外调 Reindex。
// RPC 的 rag_enabled 是无默认值的 bool，必须显式传 true——它是覆盖式写入，
// 不传等于把文档从检索里关掉。
func (l *SaveDigitalKnowledgeDocLogic) SaveDigitalKnowledgeDoc(req *types.SaveDigitalKnowledgeDocReq) (resp *types.SaveDigitalKnowledgeDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.DigitalID != "" {
		if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
			return nil, err
		}
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("标题不能为空")
	}

	if req.DocumentID == "" {
		if req.CollectionID == "" {
			return nil, errors.New("collectionId 不能为空")
		}
		content := strings.TrimSpace(req.ContentText)
		if content == "" {
			return nil, errors.New("正文不能为空")
		}
		created, err := l.svcCtx.PlatformRpc.CreateDocument(l.ctx, &platform.CreateDocumentReq{
			CollectionId: req.CollectionID,
			Title:        title,
			MediaType:    "text/plain",
			DocType:      "guide",
			ContentText:  content,
			SourceType:   "file",
			SyncMode:     "import",
			RagEnabled:   true,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "新增知识文档失败", Data: map[string]interface{}{"err": err.Error()}})
			return nil, errors.New("新增文档失败")
		}
		return &types.SaveDigitalKnowledgeDocRes{DocumentID: created.DocumentId}, nil
	}

	if _, err := l.svcCtx.PlatformRpc.UpdateDocument(l.ctx, &platform.UpdateDocumentReq{
		DocumentId:  req.DocumentID,
		Title:       title,
		ContentText: strings.TrimSpace(req.ContentText),
		RagEnabled:  true,
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "更新知识文档失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("更新文档失败")
	}
	// 「更新学习」不改正文，platform 侧只有正文变化才自动重建索引，所以显式补一次
	if req.Reindex {
		if _, err := l.svcCtx.PlatformRpc.ReindexDocument(l.ctx, &platform.ReindexDocumentReq{
			DocumentId: req.DocumentID,
		}); err != nil {
			l.logger.Error(model.LogMsg{Text: "重新索引知识文档失败", Data: map[string]interface{}{"err": err.Error()}})
			return nil, errors.New("重新学习失败")
		}
	}
	return &types.SaveDigitalKnowledgeDocRes{DocumentID: req.DocumentID}, nil
}
