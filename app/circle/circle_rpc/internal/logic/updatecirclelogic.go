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
	"fmt"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCircleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCircleLogic {
	return &UpdateCircleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCircleLogic) UpdateCircle(in *circle_rpc.UpdateCircleReq) (*circle_rpc.UpdateCircleRes, error) {
	updates := map[string]interface{}{}

	if in.IsDeleted != nil {
		updates["is_deleted"] = in.GetIsDeleted()
	}

	if len(updates) == 0 {
		return &circle_rpc.UpdateCircleRes{}, nil
	}

	updates["version"] = l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", in.CircleId)

	if err := l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", in.CircleId).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新圈子失败: %v", err)
	}

	return &circle_rpc.UpdateCircleRes{}, nil
}
