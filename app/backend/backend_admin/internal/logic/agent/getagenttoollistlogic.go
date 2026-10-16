/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package agent

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetAgentToolListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentToolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentToolListLogic {
	return &GetAgentToolListLogic{logger: beaverlog.New("get_agent_tool_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentToolListLogic) GetAgentToolList(req *types.GetAgentToolListReq) (resp *types.GetAgentToolListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAgentTools(l.ctx, &platform_rpc.ListAgentToolsReq{
		Status: int32(req.Status),
		Source: req.Source,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取工具列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	list := make([]types.AgentToolInfo, 0, len(rpcRes.Tools))
	for _, t := range rpcRes.Tools {
		list = append(list, toAgentToolInfo(t))
	}
	return &types.GetAgentToolListRes{List: list}, nil
}
