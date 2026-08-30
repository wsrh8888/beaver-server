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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListAppVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAppVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppVersionsLogic {
	return &ListAppVersionsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAppVersionsLogic) ListAppVersions(in *platform_rpc.ListAppVersionsReq) (*platform_rpc.ListAppVersionsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateArchitecture{}).Where("is_active = ?", true)
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计架构失败: %v", err)
		return nil, status.Error(codes.Internal, "获取架构总数失败")
	}

	var architectures []platform_models.UpdateArchitecture
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&architectures).Error; err != nil {
		l.Errorf("查询架构失败: %v", err)
		return nil, status.Error(codes.Internal, "获取架构列表失败")
	}

	items := make([]*platform_rpc.AppVersionsArchItem, 0, len(architectures))
	for _, arch := range architectures {
		var versions []platform_models.UpdateVersion
		if err := l.svcCtx.DB.Where("architecture_id = ?", arch.Id).Order("created_at DESC").Find(&versions).Error; err != nil {
			l.Errorf("查询架构版本失败 arch=%d: %v", arch.Id, err)
			continue
		}

		briefs := make([]*platform_rpc.AppVersionBrief, 0, len(versions))
		for _, ver := range versions {
			briefs = append(briefs, &platform_rpc.AppVersionBrief{
				VersionId: uint64(ver.Id),
				Version:   ver.Version,
			})
		}

		items = append(items, &platform_rpc.AppVersionsArchItem{
			ArchitectureId: uint64(arch.Id),
			ArchId:         uint32(arch.ArchID),
			Description:    arch.Description,
			Versions:       briefs,
		})
	}

	return &platform_rpc.ListAppVersionsRes{Total: total, Architectures: items}, nil
}
