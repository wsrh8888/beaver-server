package chat_utils

import (
	"time"

	"beaver/app/chat/chat_models"
	core "beaver/core/version"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// CreateOrUpdateConversation 更新会话信息，返回新版本号（表记录肯定存在）
func CreateOrUpdateConversation(db *gorm.DB, versionGen *core.VersionGenerator, conversationID string, conversationType int, lastSeq int64, lastMessage string) (int64, error) {
	// 生成新版本号
	version := versionGen.GetNextVersion("chat_conversation_metas", "conversation_id", conversationID)

	// 直接更新会话信息
	err := db.Model(&chat_models.ChatConversationMeta{}).
		Where("conversation_id = ?", conversationID).
		Updates(map[string]interface{}{
			"max_seq":      lastSeq,
			"last_message": lastMessage,
			"version":      version,
			"updated_at":   time.Now(),
		}).Error

	if err != nil {
		logx.Errorf("更新会话信息失败: conversationID=%s, error=%v", conversationID, err)
		return 0, err
	}

	logx.Infof("更新会话信息成功: conversationID=%s, version=%d", conversationID, version)
	return version, nil
}

// UpdateUserConversation 更新用户会话关系，返回新版本号（表记录肯定存在）
func UpdateUserConversation(db *gorm.DB, versionGen *core.VersionGenerator, userID, conversationID string, isDeleted bool) (int64, error) {
	// 生成新版本号
	version := versionGen.GetNextVersion("chat_user_conversations", "user_id", userID)

	// 准备更新字段
	updateFields := map[string]interface{}{
		"updated_at": time.Now(),
		"version":    version,
	}

	if !isDeleted {
		// 发送消息时自动恢复显示（取消隐藏状态）
		updateFields["is_hidden"] = false
	}

	// 直接更新用户会话关系
	err := db.Model(&chat_models.ChatUserConversation{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(updateFields).Error

	if err != nil {
		logx.Errorf("更新用户会话关系失败: userID=%s, conversationID=%s, error=%v", userID, conversationID, err)
		return 0, err
	}

	if !isDeleted {
		logx.Infof("更新用户会话关系成功并自动恢复显示: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
	} else {
		logx.Infof("更新用户会话关系成功: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
	}

	return version, nil
}

// UpdateAllUserConversationsInChat 更新聊天中所有用户的会话关系，返回所有用户的更新信息（表记录肯定存在）
func UpdateAllUserConversationsInChat(db *gorm.DB, versionGen *core.VersionGenerator, conversationID, senderID string, messageSeq int64) ([]UserConversationUpdate, error) {
	// 查询该会话相关的所有用户会话记录
	var userConversations []chat_models.ChatUserConversation
	err := db.Where("conversation_id = ?", conversationID).Find(&userConversations).Error
	if err != nil {
		logx.Errorf("查询会话用户关系失败: conversationID=%s, error=%v", conversationID, err)
		return nil, err
	}

	var updates []UserConversationUpdate

	// 更新每条记录，每个用户独立版本号
	for _, convo := range userConversations {
		// 总是更新版本号，因为发送消息会影响所有用户的会话状态
		version := versionGen.GetNextVersion("chat_user_conversations", "user_id", convo.UserID)

		updateFields := map[string]interface{}{
			"updated_at": time.Now(),
			"version":    version,
		}

		// 如果是隐藏状态，自动恢复显示
		if convo.IsHidden {
			updateFields["is_hidden"] = false
		}

		// 特殊处理发送者：更新已读序列号
		if convo.UserID == senderID {
			updateFields["user_read_seq"] = messageSeq
			logx.Infof("为发送者更新已读序列号: userID=%s, conversationID=%s, readSeq=%d, version=%d", convo.UserID, conversationID, messageSeq, version)
		}

		// 直接批量更新所有用户的会话关系
		err = db.Model(&chat_models.ChatUserConversation{}).
			Where("conversation_id = ? AND user_id = ?", conversationID, convo.UserID).
			Updates(updateFields).Error
		if err != nil {
			logx.Errorf("更新用户会话失败: userID=%s, conversationID=%s, error=%v", convo.UserID, conversationID, err)
			// 继续处理其他用户，不要因为一个失败而中断整个流程
		} else {
			updates = append(updates, UserConversationUpdate{
				UserID:         convo.UserID,
				ConversationID: conversationID,
				Version:        version,
			})
		}
	}

	logx.Infof("更新会话用户关系成功: conversationID=%s, 更新用户数=%d", conversationID, len(updates))
	return updates, nil
}

// UserConversationUpdate 用户会话更新信息
type UserConversationUpdate struct {
	UserID         string
	ConversationID string
	Version        int64
}
