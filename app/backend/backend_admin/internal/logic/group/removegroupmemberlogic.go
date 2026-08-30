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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMemberLogic {
	return &RemoveGroupMemberLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RemoveGroupMemberLogic) RemoveGroupMember(req *types.RemoveGroupMemberReq) (resp *types.RemoveGroupMemberRes, err error) {
	if req.GroupId == "" {
		return nil, errors.New("群组ID不能为空")
	}
	if len(req.MemberIds) == 0 {
		return nil, errors.New("成员ID列表不能为空")
	}

	for _, userID := range req.MemberIds {
		if _, err = l.svcCtx.GroupRpc.RemoveGroupMember(l.ctx, &group_rpc.RemoveGroupMemberReq{
			GroupId: req.GroupId,
			UserId:  userID,
			Kick:    true,
		}); err != nil {
			l.Errorf("移除群成员失败 groupId=%s userId=%s: %v", req.GroupId, userID, err)
			return nil, err
		}
	}
	return &types.RemoveGroupMemberRes{}, nil
}
