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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type RemoveGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRemoveGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMemberLogic {
	return &RemoveGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("remove_group_member", ctx),
	}
}

func (l *RemoveGroupMemberLogic) RemoveGroupMember(in *group_rpc.RemoveGroupMemberReq) (*group_rpc.RemoveGroupMemberRes, error) {
	var member group_models.GroupMemberModel
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?", in.GroupId, in.UserId, 1).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("成员不在群内")
		}
		l.logger.Error(model.LogMsg{Text: "查询群成员失败", Data: map[string]any{"groupId": in.GroupId, "userId": in.UserId, "err": err.Error()}})
		return nil, err
	}

	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", in.GroupId)
	if memberVersion == -1 {
		return nil, errors.New("获取群成员版本号失败")
	}

	memberStatus := int8(2)
	changeType := "leave"
	if in.Kick {
		memberStatus = 3
		changeType = "kick"
	}

	if err := l.svcCtx.DB.Model(&member).Updates(map[string]interface{}{
		"status":  memberStatus,
		"version": memberVersion,
	}).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "移除群成员失败", Data: map[string]any{"groupId": in.GroupId, "userId": in.UserId, "err": err.Error()}})
		return nil, err
	}

	operatedBy := in.OperatedBy
	if operatedBy == "" {
		operatedBy = in.UserId
	}
	_ = l.svcCtx.DB.Create(&group_models.GroupMemberChangeLogModel{
		GroupID:    in.GroupId,
		UserID:     in.UserId,
		ChangeType: changeType,
		OperatedBy: operatedBy,
		ChangeTime: time.Now(),
		Version:    memberVersion,
	}).Error

	l.logger.Info(model.LogMsg{Text: "移除群成员成功", Data: map[string]interface{}{"groupId": in.GroupId, "userId": in.UserId, "kick": in.Kick, "memberVersion": memberVersion}})

	return &group_rpc.RemoveGroupMemberRes{MemberVersion: memberVersion}, nil
}
