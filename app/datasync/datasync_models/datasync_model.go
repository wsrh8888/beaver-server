package datasync_models

import (
	"beaver/common/models"
)

// SyncCursorModel 同步游标模型 - 参考OpenIM设计
type DatasyncModel struct {
	models.Model
	UserID         string `gorm:"size:64;" json:"userId"`                      // 用户ID
	DeviceID       string `gorm:"size:128;" json:"deviceId"`                   // 设备ID（多端支持）
	DataType       string `gorm:"size:32;n" json:"dataType"`                   // 数据类型：users/friends/groups/chats/conversations
	ConversationID string `gorm:"size:64" json:"conversationId"`               // 会话ID（仅聊天消息使用）
	LastSeq        int64  `gorm:"default:0" json:"lastSeq"`                    // 这个设备最后同步的序列号（消息用）或版本号（基础数据用）
	SyncStatus     string `gorm:"size:16;default:'pending'" json:"syncStatus"` // 同步状态：pending/syncing/completed/failed
}
