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

package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/agent/agent_rpc/types/agent_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncAgentConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取 Agent 会话版本摘要（增量/全量游标）
func NewGetSyncAgentConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncAgentConversationsLogic {
	return &GetSyncAgentConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_agent_conversations", ctx),
	}
}

type agentListConversationsResult struct {
	List          []agentConversationItem `json:"list"`
	ServerVersion int64                   `json:"serverVersion"`
}

type agentConversationItem struct {
	ConversationID string `json:"conversationId"`
	Title          string `json:"title"`
	Avatar         string `json:"avatar"`
	Module         string `json:"module"`
	Workspace      string `json:"workspace"`
	UpdatedAt      int64  `json:"updatedAt"`
	LatestTopicId  string `json:"latestTopicId"`
	Preview        string `json:"preview"`
	Version        int64  `json:"version"`
}

func (l *GetSyncAgentConversationsLogic) GetSyncAgentConversations(req *types.GetSyncAgentConversationsReq) (resp *types.GetSyncAgentConversationsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.logger.Error(model.LogMsg{Text: "用户ID为空"})
		return nil, errors.New("用户ID不能为空")
	}

	serverTimestamp := time.Now().UnixMilli()
	empty := &types.GetSyncAgentConversationsRes{
		ConversationVersions: []types.AgentConversationVersionItem{},
		ServerVersion:        req.Since,
		ServerTimestamp:      serverTimestamp,
	}

	reply, err := l.svcCtx.AgentRpc.ListConversations(l.ctx, &agent_rpc.ListConversationsReq{
		UserId: userId,
		Since:  req.Since,
		Module: req.Module,
		Limit:  req.Limit,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "调用 AgentRpc.ListConversations 失败", Data: map[string]any{"userId": userId, "since": req.Since, "err": err.Error()}})
		return nil, err
	}
	if reply == nil {
		return empty, nil
	}
	if reply.Code != 0 {
		l.logger.Error(model.LogMsg{Text: "AgentRpc.ListConversations 业务失败", Data: map[string]any{"userId": userId, "code": reply.Code, "msg": reply.Msg}})
		return nil, errors.New(reply.Msg)
	}
	if reply.ResultJson == "" {
		return empty, nil
	}

	var result agentListConversationsResult
	if err := json.Unmarshal([]byte(reply.ResultJson), &result); err != nil {
		l.logger.Error(model.LogMsg{Text: "解析 AgentRpc.ListConversations 结果失败", Data: map[string]any{"userId": userId, "err": err.Error()}})
		return nil, err
	}

	versions := make([]types.AgentConversationVersionItem, 0, len(result.List))
	for _, item := range result.List {
		ws := item.Workspace
		if ws != "local" && ws != "cloud" {
			ws = "cloud"
		}
		versions = append(versions, types.AgentConversationVersionItem{
			ConversationID: item.ConversationID,
			Version:        item.Version,
			LatestTopicId:  item.LatestTopicId,
			Title:          item.Title,
			UpdatedAt:      item.UpdatedAt,
			Preview:        item.Preview,
			Module:         item.Module,
			Workspace:      ws,
		})
	}

	serverVersion := result.ServerVersion
	if serverVersion == 0 {
		serverVersion = req.Since
	}

	l.logger.Info(model.LogMsg{
		Text: "查询 Agent 会话版本完成",
		Data: map[string]any{"userId": userId, "since": req.Since, "count": len(versions), "serverVersion": serverVersion},
	})

	return &types.GetSyncAgentConversationsRes{
		ConversationVersions: versions,
		ServerVersion:        serverVersion,
		ServerTimestamp:      serverTimestamp,
	}, nil
}
