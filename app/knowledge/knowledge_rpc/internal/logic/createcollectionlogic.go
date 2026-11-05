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

type CreateCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCollectionLogic {
	return &CreateCollectionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateCollectionLogic) CreateCollection(in *knowledge_rpc.CreateCollectionReq) (*knowledge_rpc.CreateCollectionRes, error) {
	status := int(in.Status)
	if status == 0 {
		status = 1
	}
	kind := in.Kind
	if kind == "" {
		kind = "docs"
	}
	acl := in.AclMode
	if acl == "" {
		acl = "internal"
	}
	row := knowledge_models.KnowledgeCollection{
		CollectionID: "col_" + randHex(16),
		ParentID:     in.ParentId,
		Name:         in.Name,
		Kind:         kind,
		AclMode:      acl,
		IndexConfig:  in.IndexConfig,
		Description:  in.Description,
		Status:       status,
		Sort:         int(in.Sort),
	}
	if err := l.svcCtx.DB.Create(&row).Error; err != nil {
		l.Errorf("创建知识集合失败: %v", err)
		return nil, err
	}
	return &knowledge_rpc.CreateCollectionRes{CollectionId: row.CollectionID}, nil
}
