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
	"strings"

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCollectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCollectionLogic {
	return &UpdateCollectionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateCollectionLogic) UpdateCollection(in *platform_rpc.UpdateCollectionReq) (*platform_rpc.UpdateCollectionRes, error) {
	updates := map[string]any{
		"name":        in.Name,
		"kind":        in.Kind,
		"acl_mode":    in.AclMode,
		"description": in.Description,
		"status":      in.Status,
		"sort":        in.Sort,
	}
	// index_config 是 json 列：空串非法（Error 3140），且这里整列覆盖会把已有配置冲掉。
	// 调用方只改名称等情况不传配置时，保持该列不动。
	if cfg := strings.TrimSpace(in.IndexConfig); cfg != "" {
		updates["index_config"] = cfg
	}
	res := l.svcCtx.DB.Model(&knowledge_models.KnowledgeCollection{}).
		Where("collection_id = ?", in.CollectionId).
		Updates(updates)
	if res.Error != nil {
		l.Errorf("更新知识集合失败 %s: %v", in.CollectionId, res.Error)
		return nil, res.Error
	}
	return &platform_rpc.UpdateCollectionRes{}, nil
}
