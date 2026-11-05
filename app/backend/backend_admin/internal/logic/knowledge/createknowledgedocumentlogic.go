/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package knowledge

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateKnowledgeDocumentLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateKnowledgeDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKnowledgeDocumentLogic {
	return &CreateKnowledgeDocumentLogic{logger: beaverlog.New("CreateKnowledgeDocument", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateKnowledgeDocumentLogic) CreateKnowledgeDocument(req *types.CreateKnowledgeDocumentReq) (resp *types.CreateKnowledgeDocumentRes, err error) {
	rpcRes, err := l.svcCtx.KnowledgeRpc.CreateDocument(l.ctx, &knowledge_rpc.CreateDocumentReq{
		CollectionId: req.CollectionID,
		Title:        req.Title,
		MediaType:    req.MediaType,
		DocType:      req.DocType,
		SourceRef:    req.SourceRef,
		ContentText:  req.ContentText,
		MetadataText: req.MetadataText,
		RagEnabled:   req.RagEnabled,
		OwnerUserId:  req.OwnerUserID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建知识文档失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.CreateKnowledgeDocumentRes{DocumentID: rpcRes.DocumentId}, nil
}
