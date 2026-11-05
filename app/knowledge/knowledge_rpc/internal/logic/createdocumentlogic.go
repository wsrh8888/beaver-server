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

package logic

import (
	"context"

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDocumentLogic {
	return &CreateDocumentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateDocumentLogic) CreateDocument(in *knowledge_rpc.CreateDocumentReq) (*knowledge_rpc.CreateDocumentRes, error) {
	mediaType := in.MediaType
	if mediaType == "" {
		mediaType = "text/plain"
	}
	docType := in.DocType
	if docType == "" {
		docType = "guide"
	}
	// 有正文就入队等 worker 索引；空正文视为 ready（无分段可建）。
	indexStatus := "pending"
	if trim(in.ContentText) == "" {
		indexStatus = "ready"
	}
	row := knowledge_models.KnowledgeDocument{
		DocumentID:   "doc_" + randHex(16),
		CollectionID: in.CollectionId,
		Title:        in.Title,
		MediaType:    mediaType,
		DocType:      docType,
		ByteSize:     int64(len([]byte(in.ContentText))),
		SourceRef:    in.SourceRef,
		ContentText:  in.ContentText,
		MetadataText: in.MetadataText,
		RagEnabled:   in.RagEnabled,
		OwnerUserID:  in.OwnerUserId,
		Version:      1,
		UpdatedBy:    in.OwnerUserId,
		IndexStatus:  indexStatus,
		Status:       1,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.Errorf("创建知识文档失败: %v", err)
		return nil, err
	}
	return &knowledge_rpc.CreateDocumentRes{DocumentId: row.DocumentID}, nil
}
