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

package agent

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetAgentModelListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentModelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentModelListLogic {
	return &GetAgentModelListLogic{logger: beaverlog.New("get_agent_model_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentModelListLogic) GetAgentModelList(req *types.GetAgentModelListReq) (resp *types.GetAgentModelListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAgentModels(l.ctx, &platform_rpc.ListAgentModelsReq{
		Source:   req.Source,
		UserId:   req.UserID,
		Status:   int32(req.Status),
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取模型列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	list := make([]types.AgentModelInfo, 0, len(rpcRes.Models))
	for _, item := range rpcRes.Models {
		list = append(list, toAgentModelInfo(item))
	}
	return &types.GetAgentModelListRes{Total: rpcRes.Total, List: list}, nil
}
