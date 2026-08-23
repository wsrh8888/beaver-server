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

package conversation

import (
	"errors"
	"sort"
	"strings"
)

/**
 * @description: 生成会话Id
 */
func GenerateConversation(userIds []string) (string, error) {
	if len(userIds) == 1 {
		return "private_" + userIds[0], nil
	} else if len(userIds) == 2 {
		sort.Strings(userIds)
		return "private_" + strings.Join(userIds, "_"), nil
	} else {
		return "", errors.New("userIds must have a length of 1 or 2")
	}
}

/**
 * @description: 解析会话Id
 */
func ParseConversation(conversationID string) []string {
	if strings.Contains(conversationID, "_") {
		return strings.Split(conversationID, "_")
	}
	return []string{conversationID}
}

/**
 * @description: 获取会话类型
 * @return: 1: 私聊 2: 群聊 3: 圈子
 */
func GetConversationType(conversationID string) int {
	// 优先检查前缀：circle_ 圈子，group_ 群聊，private_ 私聊
	if strings.HasPrefix(conversationID, "circle_") {
		return 3
	}
	if strings.HasPrefix(conversationID, "group_") {
		return 2
	}
	if strings.HasPrefix(conversationID, "private_") {
		return 1
	}
	// 如果没有前缀，则根据是否包含下划线判断
	// 包含下划线且不是 circle_ / group_ / private_ 前缀的，通常是私聊（userId1_userId2格式）
	if strings.Contains(conversationID, "_") {
		return 1
	}
	// 不包含下划线的，通常是群聊（直接是group的UUID）
	return 2
}

/**
 * @description: 解析会话ID并返回类型和用户IDs
 * @return: conversationType (1:私聊 2:群聊 3:圈子), userIds ([]string)
 */
func ParseConversationWithType(conversationID string) (int, []string) {
	conversationType := GetConversationType(conversationID)
	userIds := ParseConversation(conversationID)

	// 对于私聊，如果是带前缀的格式 (private_userId1_userId2)，移除前缀
	if conversationType == 1 && len(userIds) >= 3 && userIds[0] == "private" {
		userIds = userIds[1:]
	}

	// 对于群聊，如果是带前缀的格式 (group_uuid)，移除前缀
	if conversationType == 2 && len(userIds) >= 2 && userIds[0] == "group" {
		userIds = userIds[1:]
	}

	// 对于圈子，如果是带前缀的格式 (circle_uuid)，移除前缀
	if conversationType == 3 && len(userIds) >= 2 && userIds[0] == "circle" {
		userIds = userIds[1:]
	}

	return conversationType, userIds
}

/**
 * @description: 从会话ID中提取对应的目标ID（对方UID或群组ID）
 * @return: targetId (对方UID或群组UUID)
 */
func GetTargetIDByConversation(conversationID string, currentUserID string) string {
	convType, ids := ParseConversationWithType(conversationID)
	if convType == 1 { // 私聊
		if len(ids) == 2 {
			if ids[0] == currentUserID {
				return ids[1]
			}
			return ids[0]
		}
		if len(ids) == 1 {
			return ids[0]
		}
	} else if convType == 2 { // 群聊
		if len(ids) > 0 {
			return ids[0]
		}
	} else if convType == 3 { // 圈子
		if len(ids) > 0 {
			return ids[0]
		}
	}
	return ""
}
