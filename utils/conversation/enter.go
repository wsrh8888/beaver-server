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
func ParseConversation(conversationId string) []string {
	if strings.Contains(conversationId, "_") {
		return strings.Split(conversationId, "_")
	}
	return []string{conversationId}
}

/**
 * @description: 获取会话类型
 * @return: 1: 私聊 2: 群聊
 */
func GetConversationType(conversationId string) int {
	if strings.Contains(conversationId, "_") {
		return 1
	}
	return 2
}
