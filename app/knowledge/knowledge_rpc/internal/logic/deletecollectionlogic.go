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

type DeleteCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCollectionLogic {
	return &DeleteCollectionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// DeleteCollection 级联：软删其下全部文档、清 OpenSearch 分段与 chunk 行，再删集合本身。
func (l *DeleteCollectionLogic) DeleteCollection(in *knowledge_rpc.DeleteCollectionReq) (*knowledge_rpc.DeleteCollectionRes, error) {
	var docIDs []string
	l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("collection_id = ? AND deleted_at = ?", in.CollectionId, zeroDeletedAt).
		Pluck("document_id", &docIDs)
	clearDocIndex(l.ctx, l.svcCtx, docIDs)

	if err := l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("collection_id = ?", in.CollectionId).
		Updates(map[string]any{"deleted_at": nowPtr(), "status": 2}).Error; err != nil {
		l.Errorf("软删集合文档失败 %s: %v", in.CollectionId, err)
		return nil, err
	}
	if err := l.svcCtx.DB.Where("collection_id = ?", in.CollectionId).
		Delete(&knowledge_models.KnowledgeCollection{}).Error; err != nil {
		l.Errorf("删除集合失败 %s: %v", in.CollectionId, err)
		return nil, err
	}
	return &knowledge_rpc.DeleteCollectionRes{}, nil
}
