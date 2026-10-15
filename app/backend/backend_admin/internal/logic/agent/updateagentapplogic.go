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

type UpdateAgentAppLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAgentAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentAppLogic {
	return &UpdateAgentAppLogic{logger: beaverlog.New("update_agent_app", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAgentAppLogic) UpdateAgentApp(req *types.UpdateAgentAppReq) (resp *types.UpdateAgentAppRes, err error) {
	_, err = l.svcCtx.PlatformRpc.UpdateAgentApp(l.ctx, &platform_rpc.UpdateAgentAppReq{
		AppId:            req.AppID,
		Name:             req.Name,
		Scene:            req.Scene,
		SystemPrompt:     req.SystemPrompt,
		LlmModelId:       req.LLMModelID,
		ModelParams:      req.ModelParams,
		Tools:            req.Tools,
		McpServerIds:     req.McpServerIDs,
		KnowledgeIds:     req.KnowledgeIDs,
		OpeningMessage:   req.OpeningMessage,
		OpeningQuestions: req.OpeningQuestions,
		Variables:        req.Variables,
		Status:           int32(req.Status),
		Sort:             int32(req.Sort),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "更新AI应用失败", Data: map[string]any{"appId": req.AppID, "err": err.Error()}})
		return nil, err
	}
	return &types.UpdateAgentAppRes{}, nil
}
