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

type CreateKnowledgeCollectionLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateKnowledgeCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateKnowledgeCollectionLogic {
	return &CreateKnowledgeCollectionLogic{logger: beaverlog.New("CreateKnowledgeCollection", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateKnowledgeCollectionLogic) CreateKnowledgeCollection(req *types.CreateKnowledgeCollectionReq) (resp *types.CreateKnowledgeCollectionRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateCollection(l.ctx, &platform_rpc.CreateCollectionReq{
		ParentId:    req.ParentID,
		Name:        req.Name,
		Kind:        req.Kind,
		AclMode:     req.AclMode,
		IndexConfig: req.IndexConfig,
		Description: req.Description,
		Status:      int32(req.Status),
		Sort:        int32(req.Sort),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建知识集合失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.CreateKnowledgeCollectionRes{CollectionID: rpcRes.CollectionId}, nil
}
