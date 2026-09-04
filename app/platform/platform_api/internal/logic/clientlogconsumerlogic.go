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
	"fmt"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/common/const/mqwsconst"
	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// ClientLogConsumerLogic 消费客户端扁平日志并写入 OpenSearch 独立索引。
type ClientLogConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewClientLogConsumerLogic(svcCtx *svc.ServiceContext) *ClientLogConsumerLogic {
	return &ClientLogConsumerLogic{
		svcCtx: svcCtx,
		logger: beaverlog.New("client_log_consumer"),
	}
}

// StartConsumer 订阅 beaver_logs，body 原样写入 OpenSearch 索引 beaver-logs。
func (l *ClientLogConsumerLogic) StartConsumer() error {
	if l.svcCtx.RocketMQ == nil {
		return fmt.Errorf("RocketMQ 未初始化")
	}
	if l.svcCtx.OpenSearch == nil {
		return fmt.Errorf("OpenSearch 未初始化")
	}
	index := l.svcCtx.Config.OpenSearch.ClientLogIndex
	if index == "" {
		index = "beaver-logs"
	}

	return l.svcCtx.RocketMQ.RegisterRawConsumer(
		mqwsconst.MqGroupClientLog,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicClientLog,
		func(body []byte) error {
			if err := l.svcCtx.OpenSearch.IndexRaw(context.Background(), index, body); err != nil {
				l.logger.Error(model.LogMsg{Text: "写入OpenSearch失败", Data: map[string]any{"err": err.Error()}})
				return err
			}
			return nil
		},
	)
}
