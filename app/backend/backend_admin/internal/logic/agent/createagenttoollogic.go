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

type CreateAgentToolLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentToolLogic {
	return &CreateAgentToolLogic{logger: beaverlog.New("create_agent_tool", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateAgentToolLogic) CreateAgentTool(req *types.CreateAgentToolReq) (resp *types.CreateAgentToolRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateAgentTool(l.ctx, &platform_rpc.CreateAgentToolReq{
		Name:     req.Name,
		ToolType: req.ToolType,
		Config:   req.Config,
		Status:   int32(req.Status),
		Sort:     int32(req.Sort),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建工具失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.CreateAgentToolRes{ToolID: rpcRes.ToolId}, nil
}
