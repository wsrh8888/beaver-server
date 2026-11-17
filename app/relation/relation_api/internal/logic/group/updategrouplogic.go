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

	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	relationtypes "beaver/app/relation/relation_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
)

type UpdateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：更新群组信息
func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateGroup", ctx),
	}
}

// UpdateGroup 只做协议转换，业务规则全部留在 relation_rpc。
func (l *UpdateGroupLogic) UpdateGroup(req *types.UpdateGroupReq) (resp *types.UpdateGroupRes, err error) {
	rpcReq := &relationtypes.UpdateGroupReq{Id: req.ID, Title: req.Title, Avatar: req.Avatar, Notice: req.Notice, Status: req.Status, MuteAll: req.MuteAll}
	_, err = l.svcCtx.GroupRpc.UpdateGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.UpdateGroupRes{}, nil
}
