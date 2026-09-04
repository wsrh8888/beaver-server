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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
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

// Track 客户端日志网关：校验量级/单条体积 → 透传至 OpenSearch（不写库）。
// 令牌校验与限流交由前置网关处理。每条日志原文（map）直接透传，
// 绝不信任客户端自报的 source/module（防伪造 source: auth_api 污染日志）。
// 契约对齐阿里云 SLS 匿名上报：logs 是自由字段数组，无桶/无多余包装。
func (l *TrackLogic) Track(req *types.TrackReq) (*types.TrackRes, error) {
	if len(req.Logs) > maxClientLogItems {
		return nil, fmt.Errorf("too many log items: %d > %d", len(req.Logs), maxClientLogItems)
	}

	clientLogger := beaverlog.New("client_log", l.ctx)

	for _, item := range req.Logs {
		b, err := json.Marshal(item)
		if err != nil {
			continue
		}
		if len(b) > maxClientLogItemLen {
			return nil, fmt.Errorf("client log item too large")
		}
		// 客户端内容一律当不可信透传
		msg := model.LogMsg{Text: "客户端日志上报", Data: item}
		switch levelOf(item) {
		case "error":
			clientLogger.Error(msg)
		case "warn":
			clientLogger.Warn(msg)
		default:
			clientLogger.Info(msg)
		}
	}

	return &types.TrackRes{}, nil
}

// levelOf 根据日志内 level 字段映射 beaverlog 严重级别，缺省 info
func levelOf(item map[string]any) string {
	switch fmt.Sprintf("%v", item["level"]) {
	case "error", "fatal":
		return "error"
	case "warn", "warning":
		return "warn"
	default:
		return "info"
	}
}
