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

type GetAgentCredentialListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetAgentCredentialListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentCredentialListLogic {
	return &GetAgentCredentialListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_agent_credential_list", ctx),
	}
}

// GetAgentCredentialList 后台只列平台官方通道（scope=platform）。个人通道由客户端自助管理，不在此暴露。
func (l *GetAgentCredentialListLogic) GetAgentCredentialList(req *types.GetAgentCredentialListReq) (resp *types.GetAgentCredentialListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAgentCredentials(l.ctx, &platform_rpc.ListAgentCredentialsReq{
		Scope:  "platform",
		Status: int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取接入通道列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	list := make([]types.AgentCredentialInfo, 0, len(rpcRes.Credentials))
	for _, item := range rpcRes.Credentials {
		list = append(list, toAgentCredentialInfo(item))
	}
	return &types.GetAgentCredentialListRes{List: list}, nil
}
