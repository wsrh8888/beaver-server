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

package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCircleLogic {
	return &DeleteCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCircleLogic) DeleteCircle(req *types.DeleteCircleReq) (resp *types.DeleteCircleRes, err error) {
	// 只有圈主可以解散
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if member.Role != 1 {
		return nil, fmt.Errorf("仅圈主可解散圈子")
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	if err = l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"version":    circleVersion,
		}).Error; err != nil {
		return nil, fmt.Errorf("解散圈子失败: %v", err)
	}

	return &types.DeleteCircleRes{}, nil
}
