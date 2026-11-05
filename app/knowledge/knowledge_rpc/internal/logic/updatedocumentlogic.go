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

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDocumentLogic {
	return &UpdateDocumentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateDocumentLogic) UpdateDocument(in *knowledge_rpc.UpdateDocumentReq) (*knowledge_rpc.UpdateDocumentRes, error) {
	var cur knowledge_models.KnowledgeDocument
	err := l.svcCtx.DB.Where("document_id = ?", in.DocumentId).First(&cur).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &knowledge_rpc.UpdateDocumentRes{}, nil
	}
	if err != nil {
		l.Errorf("查询知识文档失败 %s: %v", in.DocumentId, err)
		return nil, err
	}

	updates := map[string]any{}
	if trim(in.Title) != "" {
		updates["title"] = in.Title
	}
	if trim(in.DocType) != "" {
		updates["doc_type"] = in.DocType
	}
	if trim(in.MetadataText) != "" {
		updates["metadata_text"] = in.MetadataText
	}
	if trim(in.UpdatedBy) != "" {
		updates["updated_by"] = in.UpdatedBy
	}
	updates["rag_enabled"] = in.RagEnabled
	if in.Status != 0 {
		updates["status"] = in.Status
	}

	contentChanged := trim(in.ContentText) != "" && in.ContentText != cur.ContentText
	if contentChanged {
		updates["content_text"] = in.ContentText
		updates["byte_size"] = int64(len([]byte(in.ContentText)))
		updates["version"] = cur.Version + 1
	}

	// 是否仍参与检索：rag 开且启用（status!=2）。关闭检索/停用 → 立即清索引，不等 worker。
	stillEnabled := in.RagEnabled && (in.Status == 0 || in.Status == 1)
	if !stillEnabled {
		clearDocIndex(l.ctx, l.svcCtx, []string{in.DocumentId})
		updates["index_status"] = "ready"
	} else if contentChanged {
		updates["index_status"] = "pending"
		updates["index_next_at"] = nil
		updates["index_error"] = ""
	}

	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
			Where("document_id = ?", in.DocumentId).Updates(updates).Error; err != nil {
			l.Errorf("更新知识文档失败 %s: %v", in.DocumentId, err)
			return nil, err
		}
	}
	return &knowledge_rpc.UpdateDocumentRes{}, nil
}
