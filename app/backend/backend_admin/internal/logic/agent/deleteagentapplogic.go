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

type DeleteAgentAppLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentAppLogic {
	return &DeleteAgentAppLogic{logger: beaverlog.New("delete_agent_app", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAgentAppLogic) DeleteAgentApp(req *types.DeleteAgentAppReq) (resp *types.DeleteAgentAppRes, err error) {
	_, err = l.svcCtx.PlatformRpc.DeleteAgentApp(l.ctx, &platform_rpc.DeleteAgentAppReq{
		AppId: req.AppID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "删除AI应用失败", Data: map[string]any{"appId": req.AppID, "err": err.Error()}})
		return nil, err
	}
	return &types.DeleteAgentAppRes{}, nil
}
