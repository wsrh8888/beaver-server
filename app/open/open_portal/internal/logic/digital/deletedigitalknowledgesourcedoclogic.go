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

type DeleteDigitalKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalKnowledgeSourceDocLogic {
	return &DeleteDigitalKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital_knowledge_source_doc", ctx)}
}

// DeleteDigitalKnowledgeSourceDoc 只删知识源集合里的一条在线文档引用（软删 + 清索引），
// 不动 doc 域的在线文档本体。鉴权与归属仍按知识源行校验，避免越权删别人的引用。
func (l *DeleteDigitalKnowledgeSourceDocLogic) DeleteDigitalKnowledgeSourceDoc(req *types.DeleteDigitalKnowledgeSourceDocReq) (resp *types.DeleteDigitalKnowledgeSourceDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.DocumentID == "" {
		return nil, errors.New("documentId 不能为空")
	}
	if _, err := loadOwnedKnowledgeSource(l.svcCtx.DB, req.UserID, req.DigitalID, req.KnowledgeSourceID); err != nil {
		return nil, err
	}
	if _, rpcErr := l.svcCtx.PlatformRpc.DeleteDocument(l.ctx, &platform.DeleteDocumentReq{
		DocumentId: req.DocumentID,
	}); rpcErr != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识源在线文档引用失败", Data: map[string]interface{}{"documentId": req.DocumentID, "err": rpcErr.Error()}})
		return nil, errors.New("删除文档引用失败")
	}
	return &types.DeleteDigitalKnowledgeSourceDocRes{}, nil
}
