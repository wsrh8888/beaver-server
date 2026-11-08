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
	"errors"
	"strings"
	"time"

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/ingest"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

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

func (l *CreateDocumentLogic) CreateDocument(in *platform_rpc.CreateDocumentReq) (*platform_rpc.CreateDocumentRes, error) {
	mediaType := in.MediaType
	if mediaType == "" {
		mediaType = "text/plain"
	}
	docType := in.DocType
	if docType == "" {
		docType = "guide"
	}
	sourceDocID := trim(in.SourceDocId)
	sourceType := strings.ToLower(trim(in.SourceType))
	if sourceType == "" {
		if sourceDocID != "" {
			sourceType = "doc"
		} else {
			sourceType = "file"
		}
	}
	syncMode := strings.ToLower(trim(in.SyncMode))
	if syncMode == "" || syncMode == "direct" && sourceType != "doc" {
		syncMode = "import"
	}
	syncInterval := int(in.SyncInterval)
	if syncMode == "direct" && syncInterval <= 0 {
		syncInterval = 900
	}

	title := in.Title
	content := in.ContentText
	revision := ""
	var syncedAt *time.Time
	switch sourceType {
	case "doc":
		if sourceDocID == "" {
			return nil, errors.New("source_type=doc 时 source_doc_id 不能为空")
		}
		// 在线文档由服务端按 doc_id 读取权威快照，调用方无需再推正文。
		res, err := ingest.Fetch(l.ctx, l.svcCtx.DocRpc, ingest.Source{Type: sourceType, DocID: sourceDocID})
		if err != nil {
			l.Errorf("读取在线文档失败 doc_id=%s: %v", sourceDocID, err)
			return nil, err
		}
		content = res.Content
		revision = res.Revision
		if trim(title) == "" {
			title = res.Title
		}
		now := time.Now()
		syncedAt = &now
	case "file":
		if trim(content) == "" {
			return nil, errors.New("source_type=file 时 content_text 不能为空")
		}
	default:
		return nil, errors.New("source_type 只支持 doc|file")
	}
	if trim(title) == "" {
		return nil, errors.New("title 不能为空")
	}
	// 有正文就入队等 worker 索引；空正文视为 ready（无分段可建）。
	indexStatus := "pending"
	if trim(content) == "" {
		indexStatus = "ready"
	}
	row := knowledge_models.KnowledgeDocument{
		DocumentID:   "doc_" + randHex(16),
		CollectionID: in.CollectionId,
		Title:        title,
		MediaType:    mediaType,
		DocType:      docType,
		ByteSize:     int64(len([]byte(content))),
		SourceType:   sourceType,
		SyncMode:     syncMode,
		SourceDetail: knowledge_models.KnowledgeSourceDetail{
			DocID:        sourceDocID,
			Revision:     revision,
			SyncInterval: syncInterval,
			LastSyncedAt: syncedAt,
		},
		ContentText:   content,
		ContentDigest: ingest.Digest(content),
		MetadataText:  in.MetadataText,
		RagEnabled:    in.RagEnabled,
		IndexStatus:   indexStatus,
		Status:        1,
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.Errorf("创建知识文档失败: %v", err)
		return nil, err
	}
	return &platform_rpc.CreateDocumentRes{DocumentId: row.DocumentID}, nil
}
