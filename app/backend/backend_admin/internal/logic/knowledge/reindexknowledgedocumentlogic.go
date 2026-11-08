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
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ReindexKnowledgeDocumentLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReindexKnowledgeDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexKnowledgeDocumentLogic {
	return &ReindexKnowledgeDocumentLogic{logger: beaverlog.New("ReindexKnowledgeDocument", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ReindexKnowledgeDocumentLogic) ReindexKnowledgeDocument(req *types.ReindexKnowledgeDocumentReq) (resp *types.ReindexKnowledgeDocumentRes, err error) {
	_, err = l.svcCtx.PlatformRpc.ReindexDocument(l.ctx, &platform_rpc.ReindexDocumentReq{
		DocumentId: req.DocumentID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "重建知识文档索引失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.ReindexKnowledgeDocumentRes{}, nil
}
