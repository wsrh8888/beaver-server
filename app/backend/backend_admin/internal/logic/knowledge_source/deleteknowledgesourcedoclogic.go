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

type DeleteKnowledgeSourceDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteKnowledgeSourceDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteKnowledgeSourceDocLogic {
	return &DeleteKnowledgeSourceDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("DeleteKnowledgeSourceDoc", ctx)}
}

// DeleteKnowledgeSourceDoc 只删知识源集合里的一条文档引用（软删 + 清索引），
// 不动 doc 域的在线文档本体。归属仍按知识源行校验，避免越权删别的场景的引用。
func (l *DeleteKnowledgeSourceDocLogic) DeleteKnowledgeSourceDoc(req *types.DeleteKnowledgeSourceDocReq) (resp *types.DeleteKnowledgeSourceDocRes, err error) {
	scene, err := requireScene(req.Scene)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.DocumentID) == "" {
		return nil, errors.New("documentId 不能为空")
	}
	if _, err := loadSource(l.svcCtx.DB, scene, req.KnowledgeSourceID); err != nil {
		return nil, err
	}
	if _, rpcErr := l.svcCtx.PlatformRpc.DeleteDocument(l.ctx, &platform.DeleteDocumentReq{
		DocumentId: req.DocumentID,
	}); rpcErr != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识源文档引用失败", Data: map[string]interface{}{"documentId": req.DocumentID, "err": rpcErr.Error()}})
		return nil, errors.New("删除文档引用失败")
	}
	return &types.DeleteKnowledgeSourceDocRes{}, nil
}
