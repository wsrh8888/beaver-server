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

	"beaver/app/agent/agent_rpc/types/agent_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateModelReq      = agent_rpc.CreateModelReq
	DeleteModelReq      = agent_rpc.DeleteModelReq
	ListMessagesReq     = agent_rpc.ListMessagesReq
	ListModelsReq       = agent_rpc.ListModelsReq
	Reply               = agent_rpc.Reply
	SendMessageReq      = agent_rpc.SendMessageReq
	SubmitHostResultReq = agent_rpc.SubmitHostResultReq
	UpdateModelReq      = agent_rpc.UpdateModelReq

	Agent interface {
		SendMessage(ctx context.Context, in *SendMessageReq, opts ...grpc.CallOption) (*Reply, error)
		ListMessages(ctx context.Context, in *ListMessagesReq, opts ...grpc.CallOption) (*Reply, error)
		SubmitHostResult(ctx context.Context, in *SubmitHostResultReq, opts ...grpc.CallOption) (*Reply, error)
		CreateModel(ctx context.Context, in *CreateModelReq, opts ...grpc.CallOption) (*Reply, error)
		UpdateModel(ctx context.Context, in *UpdateModelReq, opts ...grpc.CallOption) (*Reply, error)
		DeleteModel(ctx context.Context, in *DeleteModelReq, opts ...grpc.CallOption) (*Reply, error)
		ListModels(ctx context.Context, in *ListModelsReq, opts ...grpc.CallOption) (*Reply, error)
	}

	defaultAgent struct {
		cli zrpc.Client
	}
)

func NewAgent(cli zrpc.Client) Agent {
	return &defaultAgent{cli: cli}
}

func (m *defaultAgent) SendMessage(ctx context.Context, in *SendMessageReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).SendMessage(ctx, in, opts...)
}

func (m *defaultAgent) ListMessages(ctx context.Context, in *ListMessagesReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).ListMessages(ctx, in, opts...)
}

func (m *defaultAgent) SubmitHostResult(ctx context.Context, in *SubmitHostResultReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).SubmitHostResult(ctx, in, opts...)
}

func (m *defaultAgent) CreateModel(ctx context.Context, in *CreateModelReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).CreateModel(ctx, in, opts...)
}

func (m *defaultAgent) UpdateModel(ctx context.Context, in *UpdateModelReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).UpdateModel(ctx, in, opts...)
}

func (m *defaultAgent) DeleteModel(ctx context.Context, in *DeleteModelReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).DeleteModel(ctx, in, opts...)
}

func (m *defaultAgent) ListModels(ctx context.Context, in *ListModelsReq, opts ...grpc.CallOption) (*Reply, error) {
	return agent_rpc.NewAgentClient(m.cli.Conn()).ListModels(ctx, in, opts...)
}
