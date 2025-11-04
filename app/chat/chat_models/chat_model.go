package chat_models

import (
	"beaver/common/models"
	"beaver/common/models/ctype"
	"fmt"
)

type ChatMessage struct {
	models.Model
	MessageID        string        `gorm:"size:64;uniqueIndex" json:"messageId"`           // 唯一消息ID（客户端生成+服务端确认）
	ConversationID   string        `gorm:"size:128;index" json:"conversationId"`           // 所属会话ID
	ConversationType int           `gorm:"not" json:"conversationType"`                    // 会话类型（1=私聊 2=群聊）
	Seq              int64         `gorm:"not;default:0;index" json:"seq"`                 // 消息在会话内的序列号（单调递增）
	SendUserID       *string       `gorm:"size:64;index" json:"sendUserId,omitempty"`      // 发送者用户ID（系统消息可为null）
	MsgType          ctype.MsgType `gorm:"not" json:"msgType"`                             // 消息类型（TEXT/IMAGE/VIDEO/REVOKE/DELETE/EDIT等）
	TargetMessageID  string        `gorm:"size:64;index" json:"targetMessageId,omitempty"` // 针对的原消息ID（撤回/删除/编辑事件）
	MsgPreview       string        `gorm:"size:256" json:"msgPreview"`                     // 消息预览文本
	Msg              *ctype.Msg    `gorm:"type:json" json:"msg"`                           // 消息内容（JSON）
}

func (chat ChatMessage) MsgPreviewMethod() string {
	fmt.Println("chat.Msg.Type", chat.Msg.Type)

	switch chat.Msg.Type {
	case 1:
		return chat.Msg.TextMsg.Content
	case 2:
		return "[图片消息]"
	case 3:
		return "[视频消息]"
	case 4:
		return "[文件消息]"
	case 5:
		return "[语音消息]"
	case 6:
		return "[表情消息]"
	case 7:
		return "[系统提示]"
	default:
		return "[未知消息]"
	}
}
