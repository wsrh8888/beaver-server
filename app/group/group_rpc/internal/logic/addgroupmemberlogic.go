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
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGroupMemberLogic) AddGroupMember(in *group_rpc.AddGroupMemberReq) (*group_rpc.AddGroupMemberRes, error) {
	if in.GroupId == "" || in.UserId == "" {
		return nil, errors.New("group_id 和 user_id 不能为空")
	}

	var group group_models.GroupModel
	if err := l.svcCtx.DB.Where("group_id = ? AND status = ?", in.GroupId, 1).First(&group).Error; err != nil {
		return nil, errors.New("群组不存在或已解散")
	}

	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", in.GroupId)
	if memberVersion == -1 {
		return nil, errors.New("获取群成员版本号失败")
	}

	now := time.Now()
	added := false

	var existing group_models.GroupMemberModel
	err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&existing).Error
	if err == nil {
		if existing.Status == 1 {
			return &group_rpc.AddGroupMemberRes{
				Added:         false,
				MemberVersion: existing.Version,
			}, nil
		}
		if err := l.svcCtx.DB.Model(&existing).Updates(map[string]interface{}{
			"status":    1,
			"role":      3,
			"join_time": now,
			"version":   memberVersion,
		}).Error; err != nil {
			return nil, errors.New("恢复群成员失败")
		}
		added = true
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		member := group_models.GroupMemberModel{
			GroupID:  in.GroupId,
			UserID:   in.UserId,
			Role:     3,
			Status:   1,
			JoinTime: now,
			Version:  memberVersion,
		}
		if err := l.svcCtx.DB.Create(&member).Error; err != nil {
			return nil, errors.New("添加群成员失败")
		}
		added = true
	} else {
		return nil, errors.New("查询群成员失败")
	}

	if added {
		operatedBy := in.OperatedBy
		if operatedBy == "" {
			operatedBy = in.UserId
		}
		_ = l.svcCtx.DB.Create(&group_models.GroupMemberChangeLogModel{
			GroupID:    in.GroupId,
			UserID:     in.UserId,
			ChangeType: "join",
			OperatedBy: operatedBy,
			ChangeTime: now,
			Version:    memberVersion,
		}).Error
	}

	return &group_rpc.AddGroupMemberRes{
		Added:         added,
		MemberVersion: memberVersion,
	}, nil
}
