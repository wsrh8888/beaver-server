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

type GetAgentAppListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentAppListLogic {
	return &GetAgentAppListLogic{logger: beaverlog.New("get_agent_app_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentAppListLogic) GetAgentAppList(req *types.GetAgentAppListReq) (resp *types.GetAgentAppListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAgentApps(l.ctx, &platform_rpc.ListAgentAppsReq{
		Scene:   req.Scene,
		Status:  int32(req.Status),
		Page:    int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取AI应用列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}

	list := make([]types.AgentAppInfo, 0, len(rpcRes.Apps))
	for _, a := range rpcRes.Apps {
		list = append(list, toAgentAppInfo(a))
	}
	return &types.GetAgentAppListRes{Total: rpcRes.Total, List: list}, nil
}
