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
			version, err := versionGen.GetNextVersion("conversations", "", "")
			if err != nil {
				logx.Errorf("获取会话版本号失败: %v", err)
				return err
			}
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
		version, err := versionGen.GetNextVersion("conversations", "", "")
		if err != nil {
			logx.Errorf("获取会话版本号失败: %v", err)
			return err
		}
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
func UpdateUserConversation(db *gorm.DB, versionGen *core.VersionGenerator, userID, conversationID string, isDeleted bool) error {
	var userConvo chat_models.ChatUserConversation
	err := db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&userConvo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，则创建
			version, err := versionGen.GetNextVersion("chat_user_conversations", "", "")
			if err != nil {
				logx.Errorf("获取用户会话版本号失败: %v", err)
				return err
			}
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
		version, err := versionGen.GetNextVersion("user_conversations", "", "")
		if err != nil {
			logx.Errorf("获取用户会话版本号失败: %v", err)
			return err
		}
		err = db.Model(&userConvo).
			Updates(map[string]interface{}{
				"updated_at": time.Now(),
				"version":    version,
			}).Error
		if err != nil {
			logx.Errorf("更新用户会话关系失败: %v", err)
			return err
		}
		logx.Infof("更新用户会话关系成功: userID=%s, conversationID=%s, version=%d", userID, conversationID, version)
	}
	return nil
}
