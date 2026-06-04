package openevent

import "fmt"

// 智能机器人 Robot 平台 Webhook 事件类型
const (
	EventIMMessageReceive       = "im.message.receive"
	EventIMMessageReceiveGroup  = "im.message.receive.group_at"
	EventIMBotFollowed          = "im.bot.followed"
	EventIMBotUnfollowed        = "im.bot.unfollowed"
	EventIMChatMemberBotAdded   = "im.chat.member.bot.added"
	EventIMChatMemberBotRemoved = "im.chat.member.bot.removed"
)

func ValidateRobotEventType(eventType string) error {
	switch eventType {
	case EventIMMessageReceive, EventIMMessageReceiveGroup,
		EventIMBotFollowed, EventIMBotUnfollowed,
		EventIMChatMemberBotAdded, EventIMChatMemberBotRemoved:
		return nil
	default:
		return fmt.Errorf("不支持的事件类型: %s", eventType)
	}
}
