package chat_models

import (
	"beaver/common/models"
	"beaver/common/models/ctype"
)

type ChatMessage struct {
	models.Model
	MessageID        string        `gorm:"size:64;uniqueIndex" json:"messageId"`      // 唯一消息ID（客户端生成+服务端确认）
	ConversationID   string        `gorm:"size:128;index" json:"conversationId"`      // 所属会话ID
	ConversationType int           `gorm:"not" json:"conversationType"`               // 会话类型（1=私聊 2=群聊）
	Seq              int64         `gorm:"not;default:0;index" json:"seq"`            // 消息在会话内的序列号（基于ConversationID递增，从1开始）
	SendUserID       *string       `gorm:"size:64;index" json:"sendUserId,omitempty"` // 发送者用户ID（通知消息可为null）
	MsgType          ctype.MsgType `gorm:"not" json:"msgType"`                        // 消息类型（TEXT/IMAGE/VIDEO/REVOKE/DELETE/EDIT等）
	MsgPreview       string        `gorm:"size:256" json:"msgPreview"`                // 消息预览文本
	Msg              *ctype.Msg    `gorm:"type:json" json:"msg"`                      // 消息内容（JSON）
}

func (chat ChatMessage) MsgPreviewMethod() string {
	if chat.Msg == nil {
		return "[消息]"
	}
	return getPreview(chat.Msg)
}

// getPreview 递归获取消息预览
func getPreview(msg *ctype.Msg) string {
	if msg == nil {
		return ""
	}

	switch msg.Type {
	case ctype.TextMsgType:
		if msg.TextMsg != nil {
			return msg.TextMsg.Content
		}
	case ctype.ImageMsgType:
		return "[图片消息]"
	case ctype.VideoMsgType:
		return "[视频消息]"
	case ctype.FileMsgType:
		return "[文件消息]"
	case ctype.VoiceMsgType:
		return "[语音消息]"
	case ctype.EmojiMsgType:
		return "[表情消息]"
	case ctype.NotificationMsgType:
		return "[通知消息]"
	case ctype.AudioFileMsgType:
		return "[音频文件]"
	case ctype.CallMsgType:
		return "[音视频通话]"
	case ctype.WithdrawMsgType:
		return "[撤回消息]"
	case ctype.ReplyMsgType:
		if msg.ReplyMsg != nil {
			return getPreview(msg.ReplyMsg.ReplyMsg)
		}
	case ctype.ForwardMsgType:
		return "[聊天记录]"
	}
	return "[未知消息]"
}
