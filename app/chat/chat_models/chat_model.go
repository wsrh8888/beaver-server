package chat_models

import (
	"beaver/common/models"
	"beaver/common/models/ctype"
	"errors"
	"unicode/utf8"
)

const MaxTextMessageRunes = 5000

type ChatMessage struct {
	models.Model
	MessageID        string        `gorm:"size:64;uniqueIndex" json:"messageId"`      // 唯一消息ID（客户端生成+服务端确认）
	ConversationID   string        `gorm:"size:128;index" json:"conversationId"`      // 所属会话ID
	ConversationType int           `gorm:"not null" json:"conversationType"`          // 会话类型（1=私聊 2=群聊 3=圈子）
	Seq              int64         `gorm:"not null;default:0;index" json:"seq"`       // 消息在会话内的序列号（基于ConversationID递增，从1开始）
	SendUserID       *string       `gorm:"size:64;index" json:"sendUserId,omitempty"` // 发送者用户ID（通知消息可为null）
	MsgType          ctype.MsgType `gorm:"not null" json:"msgType"`                   // 消息类型
	Status           uint8         `gorm:"not null;default:1" json:"status"`          // 消息状态：1=正常 2=已撤回 3=已编辑
	MsgPreview       string        `gorm:"size:200" json:"msgPreview"`                // 消息预览文本（最多200字符）
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
		if msg.FileMsg != nil && msg.FileMsg.FileName != "" {
			return "[文件] " + msg.FileMsg.FileName
		}
		return "[文件]"
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
	case ctype.MarkdownMsgType:
		if msg.MarkdownMsg != nil {
			if msg.MarkdownMsg.Title != "" {
				return msg.MarkdownMsg.Title
			}
			return truncatePreview(msg.MarkdownMsg.Content, 80)
		}
	case ctype.LinkMsgType:
		if msg.LinkMsg != nil {
			return "[链接] " + msg.LinkMsg.Title
		}
	case ctype.CloudDocMsgType:
		if msg.CloudDocMsg != nil {
			prefix := cloudDocPreviewLabel(msg.CloudDocMsg.DocType)
			if msg.CloudDocMsg.Title != "" {
				return prefix + msg.CloudDocMsg.Title
			}
			return prefix + "云文档"
		}
		return "[云文档]"
	}
	return "[未知消息]"
}

func cloudDocPreviewLabel(docType int) string {
	switch docType {
	case ctype.CloudDocTypeSheet:
		return "[表格] "
	case ctype.CloudDocTypeSlide:
		return "[幻灯片] "
	case ctype.CloudDocTypeMind:
		return "[思维笔记] "
	default:
		return "[文档] "
	}
}

func truncatePreview(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "…"
}

func ValidateMsgLength(msg *ctype.Msg) error {
	if msg == nil {
		return nil
	}

	switch msg.Type {
	case ctype.TextMsgType:
		if msg.TextMsg != nil && utf8.RuneCountInString(msg.TextMsg.Content) > MaxTextMessageRunes {
			return errors.New("消息内容过长")
		}
	case ctype.MarkdownMsgType:
		if msg.MarkdownMsg != nil && utf8.RuneCountInString(msg.MarkdownMsg.Content) > MaxTextMessageRunes {
			return errors.New("消息内容过长")
		}
	case ctype.ReplyMsgType:
		if msg.ReplyMsg != nil {
			if err := ValidateMsgLength(msg.ReplyMsg.ReplyMsg); err != nil {
				return err
			}
		}
	}

	return nil
}
