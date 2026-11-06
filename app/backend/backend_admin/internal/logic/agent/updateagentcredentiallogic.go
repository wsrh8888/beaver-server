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

type UpdateAgentCredentialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateAgentCredentialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentCredentialLogic {
	return &UpdateAgentCredentialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("update_agent_credential", ctx),
	}
}

// UpdateAgentCredential 改通道。apiKey 传空=不改，有值才让 platform 重算密文与 hint。
func (l *UpdateAgentCredentialLogic) UpdateAgentCredential(req *types.UpdateAgentCredentialReq) (resp *types.UpdateAgentCredentialRes, err error) {
	if req.CredentialID == "" {
		return nil, errors.New("通道 ID 不能为空")
	}
	_, err = l.svcCtx.PlatformRpc.UpdateAgentCredential(l.ctx, &platform_rpc.UpdateAgentCredentialReq{
		CredentialId: req.CredentialID,
		Name:         req.Name,
		ModelService: req.ModelService,
		ChannelType:  req.ChannelType,
		BaseUrl:      req.BaseURL,
		ApiKey:       req.APIKey,
		Sort:         int32(req.Sort),
		Status:       int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "更新接入通道失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.UpdateAgentCredentialRes{}, nil
}
