package chat_utils

import (
	"time"

	"beaver/app/chat/chat_models"
	core "beaver/core/version"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// CreateOrUpdateConversation 创建或更新会话信息
func CreateOrUpdateConversation(db *gorm.DB, versionGen *core.VersionGenerator, conversationID string, conversationType int, lastSeq int64, lastMessage string) error {
	var convModel chat_models.ChatConversationMeta
	err := db.Where("conversation_id = ?", conversationID).First(&convModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，则创建
			version := versionGen.GetNextVersion("chat_conversation_metas", "conversation_id", conversationID)
			convModel = chat_models.ChatConversationMeta{
				ConversationID: conversationID,
				Type:           conversationType,
				MaxSeq:         lastSeq,
				LastMessage:    lastMessage,
				Version:        version,
			}
			err = db.Create(&convModel).Error
			if err != nil {
				logx.Errorf("创建会话信息失败: %v", err)
				return err
			}
			logx.Infof("创建会话信息成功: conversationID=%s, version=%d", conversationID, version)
		} else {
			logx.Errorf("查询会话信息失败: %v", err)
			return err
		}
	} else {
		// 如果存在，则更新
		version := versionGen.GetNextVersion("chat_conversation_metas", "conversation_id", conversationID)
		err = db.Model(&convModel).
			Updates(map[string]interface{}{
				"max_seq":      lastSeq,
				"last_message": lastMessage,
				"version":      version,
				"updated_at":   time.Now(),
			}).Error
		if err != nil {
			logx.Errorf("更新会话信息失败: %v", err)
			return err
		}
		logx.Infof("更新会话信息成功: conversationID=%s, version=%d", conversationID, version)
	}
	return nil
}

// UpdateUserConversation 更新用户会话关系（不再更新LastMessage，只更新版本号）
// isDeleted: 是否删除会话（true=隐藏，false=正常/发送消息时恢复显示）
func UpdateUserConversation(db *gorm.DB, versionGen *core.VersionGenerator, userID, conversationID string, isDeleted bool) error {
	var userConvo chat_models.ChatUserConversation
	err := db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&userConvo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，则创建
			// 注意：version基于user_id递增，所有用户的会话设置共享同一个版本号（用户级同步）
			version := versionGen.GetNextVersion("chat_user_conversations", "user_id", userID)
			err = db.Create(&chat_models.ChatUserConversation{
				UserID:         userID,
				ConversationID: conversationID,
				IsHidden:       isDeleted, // 兼容旧的isDeleted参数
				IsPinned:       false,
				IsMuted:        false,
				UserReadSeq:    0,
				Version:        version,
			}).Error
			if err != nil {
				logx.Errorf("创建用户会话关系失败: %v", err)
				return err
			}
			logx.Infof("创建用户会话关系成功: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
		} else {
			logx.Errorf("查询用户会话关系失败: %v", err)
			return err
		}
	} else {
		// 如果存在，则更新
		// 注意：version基于user_id递增，所有用户的会话设置共享同一个版本号（用户级同步）
		version := versionGen.GetNextVersion("chat_user_conversations", "user_id", userID)

		// 如果不是删除操作（即发送消息等正常操作），自动恢复显示状态
		updateFields := map[string]interface{}{
			"updated_at": time.Now(),
			"version":    version,
		}

		if !isDeleted {
			// 发送消息时自动恢复显示（取消隐藏状态）
			updateFields["is_hidden"] = false
		}

		err = db.Model(&userConvo).Updates(updateFields).Error
		if err != nil {
			logx.Errorf("更新用户会话关系失败: %v", err)
			return err
		}
		if !isDeleted {
			logx.Infof("更新用户会话关系成功并自动恢复显示: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
		} else {
			logx.Infof("更新用户会话关系成功: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
		}
	}
	return nil
}

// UpdateAllUserConversationsInChat 更新聊天中所有用户的会话关系（发送消息时自动恢复隐藏状态）
func UpdateAllUserConversationsInChat(db *gorm.DB, versionGen *core.VersionGenerator, conversationID string, excludeUserID string) error {
	// 查询需要更新的记录
	var userConversations []chat_models.ChatUserConversation
	err := db.Where("conversation_id = ? AND user_id != ? AND is_hidden = ?",
		conversationID, excludeUserID, true).Find(&userConversations).Error
	if err != nil {
		logx.Errorf("查询需要更新的用户会话失败: conversationID=%s, error=%v", conversationID, err)
		return err
	}

	if len(userConversations) == 0 {
		logx.Infof("没有需要恢复的隐藏会话: conversationID=%s", conversationID)
		return nil
	}

	// 直接逐个更新每条记录（N次数据库操作）
	for _, convo := range userConversations {
		version := versionGen.GetNextVersion("chat_user_conversations", "user_id", convo.UserID)
		err = db.Model(&chat_models.ChatUserConversation{}).
			Where("conversation_id = ? AND user_id = ?", conversationID, convo.UserID).
			Updates(map[string]interface{}{
				"is_hidden":  false,
				"updated_at": time.Now(),
				"version":    version,
			}).Error
		if err != nil {
			logx.Errorf("更新用户会话失败: userID=%s, conversationID=%s, error=%v", convo.UserID, conversationID, err)
			// 继续处理其他用户，不要因为一个失败而中断整个流程
		}
	}

	logx.Infof("恢复隐藏会话成功: conversationID=%s, 恢复用户数=%d", conversationID, len(userConversations))
	return nil
}
