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
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateAgentCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateAgentCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentCredentialLogic {
	return &CreateAgentCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("create_agent_credential", ctx),
	}
}

// CreateAgentCredential 后台新建官方接入通道（scope=platform）。密钥明文转发给 platform，由其加密落库。
func (l *CreateAgentCredentialLogic) CreateAgentCredential(req *types.CreateAgentCredentialReq) (resp *types.CreateAgentCredentialRes, err error) {
	if req.Name == "" || req.BaseURL == "" {
		return nil, errors.New("通道名称与接口地址不能为空")
	}
	rpcRes, err := l.svcCtx.PlatformRpc.CreateAgentCredential(l.ctx, &platform_rpc.CreateAgentCredentialReq{
		Name:         req.Name,
		ModelService: req.ModelService,
		ChannelType:  req.ChannelType,
		BaseUrl:      req.BaseURL,
		ApiKey:       req.APIKey,
		Scope:        "platform",
		Sort:         int32(req.Sort),
		Status:       int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建接入通道失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	if rpcRes.CredentialId == "" {
		return nil, errors.New("创建接入通道失败")
	}
	return &types.CreateAgentCredentialRes{CredentialID: rpcRes.CredentialId}, nil
}
