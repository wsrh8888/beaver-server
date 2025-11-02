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
		return userIds[0], nil
	} else if len(userIds) == 2 {
		sort.Strings(userIds)
		return strings.Join(userIds, "_"), nil
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
 * @return: 1: 私聊 2: 群聊
 */
func GetConversationType(conversationID string) int {
	if strings.Contains(conversationID, "_") {
		return 1
	}
	return 2
}

/**
 * @description: 解析会话ID并返回类型和用户IDs
 * @return: conversationType (1:私聊 2:群聊), userIds ([]string)
 */
func ParseConversationWithType(conversationID string) (int, []string) {
	conversationType := GetConversationType(conversationID)
	userIds := ParseConversation(conversationID)

	// 对于私聊，如果是带前缀的格式 (private_userId1_userId2)，移除前缀
	if conversationType == 1 && len(userIds) >= 3 && userIds[0] == "private" {
		userIds = userIds[1:]
	}

	return conversationType, userIds
}
