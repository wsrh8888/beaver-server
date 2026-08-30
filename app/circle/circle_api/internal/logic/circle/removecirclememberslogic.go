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

type RemoveCircleMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCircleMembersLogic {
	return &RemoveCircleMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveCircleMembersLogic) RemoveCircleMembers(req *types.RemoveCircleMembersReq) (resp *types.RemoveCircleMembersRes, err error) {
	if len(req.UserIds) == 0 {
		return nil, fmt.Errorf("请选择要移除的成员")
	}

	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可移除成员")
	}

	var members []circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id IN ?", req.CircleID, req.UserIds).Find(&members).Error; err != nil {
		return nil, fmt.Errorf("查询成员失败")
	}

	for _, member := range members {
		if member.UserID == req.UserID {
			return nil, fmt.Errorf("不能移除自己")
		}
		if operator.Role == 1 {
			if member.Role == 1 {
				return nil, fmt.Errorf("不能移除圈主")
			}
		} else if operator.Role == 2 {
			if member.Role != 3 {
				return nil, fmt.Errorf("管理员只能移除普通成员")
			}
		}
	}

	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id IN ?", req.CircleID, req.UserIds).
		Delete(&circle_models.CircleMemberModel{}).Error; err != nil {
		return nil, fmt.Errorf("移除成员失败: %v", err)
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Update("version", circleVersion)

	return &types.RemoveCircleMembersRes{}, nil
}
