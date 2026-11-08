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

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDocumentLogic {
	return &DeleteDocumentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// DeleteDocument 软删文档 + 同步清 OpenSearch 分段与 chunk 行，检索立即不再命中。
func (l *DeleteDocumentLogic) DeleteDocument(in *platform_rpc.DeleteDocumentReq) (*platform_rpc.DeleteDocumentRes, error) {
	clearDocIndex(l.ctx, l.svcCtx, []string{in.DocumentId})
	if err := l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ?", in.DocumentId).
		Updates(map[string]any{"deleted_at": nowPtr(), "status": 2}).Error; err != nil {
		l.Errorf("删除知识文档失败 %s: %v", in.DocumentId, err)
		return nil, err
	}
	return &platform_rpc.DeleteDocumentRes{}, nil
}
