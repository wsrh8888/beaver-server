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

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ReindexDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReindexDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexDocumentLogic {
	return &ReindexDocumentLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReindexDocumentLogic) ReindexDocument(in *platform_rpc.ReindexDocumentReq) (*platform_rpc.ReindexDocumentRes, error) {
	if trim(in.DocumentId) == "" {
		return nil, errors.New("document_id 为空")
	}
	var row knowledge_models.KnowledgeDocument
	err := l.svcCtx.DB.Where("document_id = ? AND deleted_at = ?", in.DocumentId, zeroDeletedAt).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("文档不存在")
	}
	if err != nil {
		l.Errorf("查询知识文档失败 %s: %v", in.DocumentId, err)
		return nil, err
	}
	// 重新入队：清退避与错误计数，置 pending 让 worker 下一轮立即捞取。
	res := l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).
		Where("document_id = ?", in.DocumentId).
		Updates(map[string]any{
			"index_status":      "pending",
			"index_retry_count": 0,
			"index_error":       "",
			"index_next_at":     nil,
		})
	if res.Error != nil {
		l.Errorf("重置索引状态失败 %s: %v", in.DocumentId, res.Error)
		return nil, res.Error
	}
	return &platform_rpc.ReindexDocumentRes{}, nil
}
