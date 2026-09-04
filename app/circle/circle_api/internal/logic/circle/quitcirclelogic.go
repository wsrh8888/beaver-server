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

type QuitCircleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewQuitCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitCircleLogic {
	return &QuitCircleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("quit_circle", ctx),
	}
}

func (l *QuitCircleLogic) QuitCircle(req *types.QuitCircleReq) (resp *types.QuitCircleRes, err error) {
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("你不是该圈子成员")
	}
	if member.Role == 1 {
		return nil, fmt.Errorf("圈主不能退出，请先转让圈主")
	}

	if err = l.svcCtx.DB.Delete(&member).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "退出圈子失败", Data: map[string]any{"circleId": req.CircleID, "userId": req.UserID, "err": err.Error()}})
		return nil, fmt.Errorf("退出圈子失败: %v", err)
	}

	// 更新圈子版本
	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Update("version", circleVersion)

	l.logger.Info(model.LogMsg{Text: "退出圈子成功", Data: map[string]interface{}{"circleId": req.CircleID, "userId": req.UserID}})

	return &types.QuitCircleRes{}, nil
}
