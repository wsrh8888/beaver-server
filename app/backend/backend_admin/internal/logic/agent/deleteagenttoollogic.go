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

type DeleteAgentToolLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentToolLogic {
	return &DeleteAgentToolLogic{logger: beaverlog.New("delete_agent_tool", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAgentToolLogic) DeleteAgentTool(req *types.DeleteAgentToolReq) (resp *types.DeleteAgentToolRes, err error) {
	_, err = l.svcCtx.PlatformRpc.DeleteAgentTool(l.ctx, &platform_rpc.DeleteAgentToolReq{
		ToolId: req.ToolID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "删除工具失败", Data: map[string]any{"toolId": req.ToolID, "err": err.Error()}})
		return nil, err
	}
	return &types.DeleteAgentToolRes{}, nil
}
