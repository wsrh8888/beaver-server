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

type GetAgentToolLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentToolLogic {
	return &GetAgentToolLogic{logger: beaverlog.New("get_agent_tool", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentToolLogic) GetAgentTool(req *types.GetAgentToolReq) (resp *types.GetAgentToolRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetAgentTool(l.ctx, &platform_rpc.GetAgentToolReq{
		ToolId: req.ToolID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取工具失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	out := types.GetAgentToolRes{}
	if rpcRes.Tool != nil {
		out.Tool = toAgentToolInfo(rpcRes.Tool)
	}
	return &out, nil
}
