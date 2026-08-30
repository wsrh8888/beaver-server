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

type UpdateMemberRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberRoleLogic {
	return &UpdateMemberRoleLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateMemberRoleLogic) UpdateMemberRole(req *types.UpdateMemberRoleReq) (resp *types.UpdateMemberRoleRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("成员ID不能为空")
	}
	if req.Role < 1 || req.Role > 3 {
		return nil, errors.New("角色值无效，有效值为1-3")
	}

	role := int32(req.Role)
	_, err = l.svcCtx.GroupRpc.UpdateGroupMember(l.ctx, &group_rpc.UpdateGroupMemberReq{
		Id:   uint64(req.Id),
		Role: &role,
	})
	if err != nil {
		l.Errorf("更新成员角色失败: %v", err)
		return nil, err
	}
	return &types.UpdateMemberRoleRes{}, nil
}
