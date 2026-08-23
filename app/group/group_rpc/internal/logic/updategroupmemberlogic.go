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

type UpdateGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupMemberLogic {
	return &UpdateGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateGroupMemberLogic) UpdateGroupMember(in *group_rpc.UpdateGroupMemberReq) (*group_rpc.UpdateGroupMemberRes, error) {
	var member group_models.GroupMemberModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("群组成员不存在")
		}
		l.Errorf("查询群组成员失败: %v", err)
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.Role != nil {
		updates["role"] = *in.Role
	}
	if in.MuteMinutes != nil {
		var mutedUntil *time.Time
		if *in.MuteMinutes > 0 {
			t := time.Now().Add(time.Duration(*in.MuteMinutes) * time.Minute)
			mutedUntil = &t
		}
		updates["muted_until"] = mutedUntil
	}
	if len(updates) == 0 {
		return &group_rpc.UpdateGroupMemberRes{}, nil
	}

	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", member.GroupID)
	if memberVersion == -1 {
		return nil, errors.New("获取群成员版本号失败")
	}
	updates["version"] = memberVersion

	if err := l.svcCtx.DB.Model(&member).Updates(updates).Error; err != nil {
		l.Errorf("更新群成员失败: %v", err)
		return nil, err
	}

	changeType := "role_change"
	if in.MuteMinutes != nil {
		changeType = "mute"
	}
	_ = l.svcCtx.DB.Create(&group_models.GroupMemberChangeLogModel{
		GroupID:    member.GroupID,
		UserID:     member.UserID,
		ChangeType: changeType,
		OperatedBy: "system",
		ChangeTime: time.Now(),
		Version:    memberVersion,
	}).Error

	return &group_rpc.UpdateGroupMemberRes{}, nil
}
