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
	"fmt"
	"strings"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoLogic {
	return &GroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInfoLogic) GroupInfo(req *types.GroupInfoReq) (resp *types.GroupInfoRes, err error) {
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, "group_id = ?", req.GroupID).Error

	if err != nil {
		logx.Errorf("查询群组失败: %s", err.Error())
		return nil, errors.New("群组不存在")
	}

	var memberCount int64
	_ = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id = ?", req.GroupID).Count(&memberCount).Error

	resp = &types.GroupInfoRes{
		GroupID:        group.GroupID,
		Title:          group.Title,
		Avatar:         group.Avatar,
		ConversationID: group.GroupID,
		MemberCount:    int(memberCount),
		CreatorID:      group.CreatorID,
		Notice:         group.Notice,
		JoinType:       group.JoinType,
		Status:         group.Status,
		CreatedAt:      time.Time(group.CreatedAt).Unix(),
		UpdatedAt:      time.Time(group.UpdatedAt).Unix(),
		Version:        group.Version,
	}

	var member group_models.GroupMemberModel
	if l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", req.GroupID, req.UserID).First(&member).Error == nil {
		var invite group_models.GroupInviteLinkModel
		if e := l.svcCtx.DB.Where("group_id = ? AND status = 1", group.GroupID).Order("id asc").First(&invite).Error; e == nil {
			domain := strings.TrimRight(strings.TrimSpace(l.svcCtx.Config.Domain), "/")
			if domain != "" && invite.Token != "" {
				resp.InviteUrl = fmt.Sprintf("%s/api/group/v1/invite_code?code=%s", domain, invite.Token)
			}
		}
	}
	return resp, nil
}
