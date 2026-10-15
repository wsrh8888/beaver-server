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

type UpdateAgentToolLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentToolLogic {
	return &UpdateAgentToolLogic{logger: beaverlog.New("update_agent_tool", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAgentToolLogic) UpdateAgentTool(req *types.UpdateAgentToolReq) (resp *types.UpdateAgentToolRes, err error) {
	_, err = l.svcCtx.PlatformRpc.UpdateAgentTool(l.ctx, &platform_rpc.UpdateAgentToolReq{
		ToolId:   req.ToolID,
		Name:     req.Name,
		ToolType: req.ToolType,
		Config:   req.Config,
		Status:   int32(req.Status),
		Sort:     int32(req.Sort),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "更新工具失败", Data: map[string]any{"toolId": req.ToolID, "err": err.Error()}})
		return nil, err
	}
	return &types.UpdateAgentToolRes{}, nil
}
