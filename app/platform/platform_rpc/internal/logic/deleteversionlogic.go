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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DeleteVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVersionLogic {
	return &DeleteVersionLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_version", ctx)}
}

func (l *DeleteVersionLogic) DeleteVersion(in *platform_rpc.DeleteVersionReq) (*platform_rpc.DeleteVersionRes, error) {
	if in.VersionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "版本ID不能为空")
	}

	var version platform_models.UpdateVersion
	if err := l.svcCtx.DB.First(&version, in.VersionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "版本不存在")
		}
		l.logger.Error(model.LogMsg{Text: "查询版本失败", Data: map[string]any{"versionId": in.VersionId, "err": err.Error()}})
		return nil, err
	}

	var policies []platform_models.UpdateReleasePolicy
	if err := l.svcCtx.DB.Find(&policies).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询发版策略失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	vid := uint(in.VersionId)
	for _, p := range policies {
		if p.StableVersionID == vid || p.GrayVersionID == vid {
			return nil, status.Error(codes.FailedPrecondition, "版本已被发版策略引用，请先调整策略")
		}
	}

	if err := l.svcCtx.DB.Delete(&version).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除版本失败", Data: map[string]any{"versionId": in.VersionId, "err": err.Error()}})
		return nil, status.Error(codes.Internal, "删除版本失败")
	}

	l.logger.Info(model.LogMsg{Text: "删除版本成功", Data: map[string]interface{}{"versionId": in.VersionId}})

	return &platform_rpc.DeleteVersionRes{}, nil
}
