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

type DeleteDigitalKnowledgeDocLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalKnowledgeDocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalKnowledgeDocLogic {
	return &DeleteDigitalKnowledgeDocLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital_knowledge_doc", ctx)}
}

// DeleteDigitalKnowledgeDoc 删除知识文档（软删 + 同步清索引）。
func (l *DeleteDigitalKnowledgeDocLogic) DeleteDigitalKnowledgeDoc(req *types.DeleteDigitalKnowledgeDocReq) (resp *types.DeleteDigitalKnowledgeDocRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.DigitalID != "" {
		if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
			return nil, err
		}
	}
	if req.DocumentID == "" {
		return nil, errors.New("documentId 不能为空")
	}
	if _, err := l.svcCtx.PlatformRpc.DeleteDocument(l.ctx, &platform.DeleteDocumentReq{
		DocumentId: req.DocumentID,
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "删除知识文档失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("删除文档失败")
	}
	return &types.DeleteDigitalKnowledgeDocRes{}, nil
}
