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

package track_public

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/const/mqwsconst"
)

// 字段上限：限制单次请求规模与单条日志体积，防止异常 payload 撑爆 OpenSearch 索引映射。
const (
	maxClientLogItems   = 200
	maxClientLogItemLen = 16 * 1024
)

type TrackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrackLogic {
	return &TrackLogic{ctx: ctx, svcCtx: svcCtx}
}

// Track 客户端日志网关：校验量级/单条体积 → 扁平 JSON 原样写入 RocketMQ。
// 不走 beaverlog / OTel，避免 attributes 包装与服务端日志混入同一索引。
// 消费端写入索引 beaver-logs，字段顶层扁平，查询对齐 SLS。
func (l *TrackLogic) Track(req *types.TrackReq) (*types.TrackRes, error) {
	if l.svcCtx.RocketMQ == nil {
		return nil, fmt.Errorf("client log rocketmq not configured")
	}
	if len(req.Logs) > maxClientLogItems {
		return nil, fmt.Errorf("too many log items: %d > %d", len(req.Logs), maxClientLogItems)
	}

	for _, item := range req.Logs {
		if item == nil {
			continue
		}
		b, err := json.Marshal(item)
		if err != nil {
			continue
		}
		if len(b) > maxClientLogItemLen {
			return nil, fmt.Errorf("client log item too large")
		}
		// 原样透传：不加 Text/Data/attributes，不注入 source/module/traceId
		if err := l.svcCtx.RocketMQ.SendRawJSON(l.ctx, mqwsconst.MqTopicClientLog, item); err != nil {
			return nil, fmt.Errorf("forward client logs failed: %w", err)
		}
	}
	return &types.TrackRes{}, nil
}
