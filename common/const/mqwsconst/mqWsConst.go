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

package mqwsconst

// RocketMQ Topic 类型（用于 WebSocket 推送）
type TopicType string

// RocketMQ Topic 常量
const (
	// MqTopicWs WebSocket 推送专用 Topic（用户侧）
	MqTopicWs TopicType = "ws_push_topic"
	// MqTopicWsBot 机器人长连接推送专用 Topic（开放平台侧）
	// 与用户侧 ws_push_topic 分离，用户侧 ws_api 消费者完全看不到机器人事件
	MqTopicWsBot TopicType = "ws_bot_push_topic"
	// MqTopicClientLog 客户端日志扁平 JSON（无 Message 信封）
	MqTopicClientLog TopicType = "beaver_logs"
)

// RocketMQ Consumer Group 类型
type GroupType string

// RocketMQ Consumer Group 常量
const (
	// MqGroupWs WS API 消费者组
	MqGroupWs GroupType = "ws_api_consumer_group"
	// MqGroupWsBot 开放平台机器人长连接消费者组
	MqGroupWsBot GroupType = "open_api_bot_consumer_group"
	// MqGroupClientLog 客户端日志写入 OpenSearch 的消费者组
	MqGroupClientLog GroupType = "platform_client_log_group"
)
