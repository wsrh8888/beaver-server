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

package digital

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateSystemDigitalLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSystemDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSystemDigitalLogic {
	return &CreateSystemDigitalLogic{logger: beaverlog.New("create_system_digital", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateSystemDigitalLogic) CreateSystemDigital(req *types.CreateSystemDigitalReq) (resp *types.CreateSystemDigitalRes, err error) {
	rpcRes, err := l.svcCtx.OpenRpc.CreateDigital(l.ctx, &open_rpc.CreateDigitalReq{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		Runtime:      req.Runtime,
		Visibility:   req.Visibility,
		Source:       open_models.DigitalSourceSystem,
		OwnerId:      req.UserID,
		SystemPrompt: req.SystemPrompt,
		LlmModelId:   req.LLMModelID,
		Tools:        req.Tools,
		McpServerIds: req.McpServerIDs,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建系统数字伙伴失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.CreateSystemDigitalRes{
		DigitalID: rpcRes.DigitalId,
		AppKey:    rpcRes.AppKey,
		AppSecret: rpcRes.AppSecret,
	}, nil
}
