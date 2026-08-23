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

type UpdateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateGroupLogic) UpdateGroup(req *types.UpdateGroupReq) (resp *types.UpdateGroupRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("群组ID不能为空")
	}

	rpcReq := &group_rpc.UpdateGroupReq{
		Id:     uint64(req.Id),
		Title:  req.Title,
		Avatar: req.FileName,
		Notice: req.Notice,
		Status: int32(req.Status),
	}
	muteAll := req.MuteAll
	rpcReq.MuteAll = &muteAll

	_, err = l.svcCtx.GroupRpc.UpdateGroup(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("更新群组信息失败: %v", err)
		return nil, err
	}
	return &types.UpdateGroupRes{}, nil
}
