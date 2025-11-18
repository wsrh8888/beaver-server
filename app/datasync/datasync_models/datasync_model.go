package datasync_models

import (
	"beaver/common/models"
)

// DatasyncModel 数据同步模型
type DatasyncModel struct {
	models.Model
	DataType     string `gorm:"size:32;unique" json:"dataType"`              // 数据类型：users/friends/groups/chats/conversations（唯一键）
	LastSeq      int64  `gorm:"default:0" json:"lastSeq"`                    // 最后同步的序列号（消息用）或版本号（基础数据用）
	LastSyncTime int64  `gorm:"default:0" json:"lastSyncTime"`               // 最后同步的时间戳
	SyncStatus   string `gorm:"size:16;default:'pending'" json:"syncStatus"` // 同步状态：pending/syncing/completed/failed
}
