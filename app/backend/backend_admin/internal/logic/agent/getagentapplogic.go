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

type GetAgentAppLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentAppLogic {
	return &GetAgentAppLogic{logger: beaverlog.New("get_agent_app", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentAppLogic) GetAgentApp(req *types.GetAgentAppReq) (resp *types.GetAgentAppRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetAgentApp(l.ctx, &platform_rpc.GetAgentAppReq{
		AppId: req.AppID,
		Scene: req.Scene,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取AI应用失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	out := types.GetAgentAppRes{}
	if rpcRes.App != nil {
		out.App = toAgentAppInfo(rpcRes.App)
	}
	return &out, nil
}
