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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type UpdateCircleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCircleLogic {
	return &UpdateCircleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("update_circle", ctx),
	}
}

func (l *UpdateCircleLogic) UpdateCircle(req *types.UpdateCircleReq) (resp *types.UpdateCircleRes, err error) {
	// 权限校验：必须是圈主或管理员
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if member.Role > 2 {
		return nil, fmt.Errorf("无权限，仅圈主和管理员可修改")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.JoinType != 0 {
		updates["join_type"] = req.JoinType
	}

	if len(updates) == 0 {
		return &types.UpdateCircleRes{}, nil
	}

	updates["version"] = l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	if err = l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Updates(updates).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "更新圈子失败", Data: map[string]any{"circleId": req.CircleID, "err": err.Error()}})
		return nil, fmt.Errorf("更新圈子失败: %v", err)
	}

	l.logger.Info(model.LogMsg{Text: "更新圈子成功", Data: map[string]interface{}{"circleId": req.CircleID, "userId": req.UserID}})

	return &types.UpdateCircleRes{}, nil
}
