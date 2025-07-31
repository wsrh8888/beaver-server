package chat_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
	"beaver/common/models/ctype"
	"fmt"
)

type ChatModel struct {
	models.Model
	MessageID      string                `json:"messageId"`                        // 客户端消息ID
	ConversationID string                `json:"conversationId"`                   // 会话id（单聊为用户id，群聊为群id）
	SendUserID     string                `gorm:"size:64;index"  json:"sendUserId"` // 发送者用户id
	MsgType        ctype.MsgType         `json:"msgType"`                          // 消息类型
	MsgPreview     string                `gorm:"size:64" json:"msgPreview"`        // 消息预览
	Msg            *ctype.Msg            `json:"msg"`                              // 消息内容
	SendUserModel  user_models.UserModel `gorm:"foreignKey:SendUserID;references:UUID" json:"-"`
	IsDeleted      bool                  `gorm:"not null;default:false" json:"isDeleted"` // 标记用户是否删除会话
}

func (chat ChatModel) MsgPreviewMethod() string {
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
	default:
		return "[未知消息]"
	}
}
