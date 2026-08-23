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

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetParticipantsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取房间当前成员列表
func NewGetParticipantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParticipantsLogic {
	return &GetParticipantsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetParticipantsLogic) GetParticipants(req *types.GetCallParticipantsReq) (resp *types.GetCallParticipantsRes, err error) {
	rpcResp, err := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{
		RoomId: req.RoomID,
	})
	if err != nil {
		return nil, err
	}
	participants := make([]types.Participant, 0)
	for _, p := range rpcResp.Participants {

		participants = append(participants, types.Participant{
			UserID: p.UserId,
			Status: p.Status,
		})
	}

	return &types.GetCallParticipantsRes{
		Participants: participants,
	}, nil
}
