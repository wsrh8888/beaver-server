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

type RemoveCircleMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCircleMembersLogic {
	return &RemoveCircleMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveCircleMembersLogic) RemoveCircleMembers(in *circle_rpc.RemoveCircleMembersReq) (*circle_rpc.RemoveCircleMembersRes, error) {
	if in.CircleId == "" {
		return nil, fmt.Errorf("圈子ID不能为空")
	}
	if len(in.UserIds) == 0 {
		return nil, fmt.Errorf("请选择要移除的成员")
	}

	var circle circle_models.CircleModel
	if err := l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", in.CircleId).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	for _, uid := range in.UserIds {
		if uid == circle.CreatorID {
			return nil, fmt.Errorf("不能移除圈主")
		}
	}

	if err := l.svcCtx.DB.Where("circle_id = ? AND user_id IN ?", in.CircleId, in.UserIds).
		Delete(&circle_models.CircleMemberModel{}).Error; err != nil {
		return nil, fmt.Errorf("移除成员失败: %v", err)
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", in.CircleId)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", in.CircleId).
		Update("version", circleVersion)

	return &circle_rpc.RemoveCircleMembersRes{}, nil
}
