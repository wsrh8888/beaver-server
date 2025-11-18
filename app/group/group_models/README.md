1、group_mute_list_model.go 群成员禁言列表
package group_models

import (
	"beaver/common/models"
	"time"
)

// 群成员禁言列表
type GroupMuteListModel struct {
	models.Model
	GroupID string    `gorm:"size:64;index" json:"groupId"`
	UserID  string    `gorm:"size:64;index" json:"userId"` // 被禁言的用户ID
	MutedBy string    `gorm:"size:64" json:"mutedBy"`      // 操作者
	StartAt time.Time `json:"startAt"`                     // 禁言开始时间
	EndAt   time.Time `json:"endAt"`                       // 禁言结束时间
}


2、 group_file_model 群文件
// GroupFileModel 群文件关联表：仅关联已上传的文件到群文件区
type GroupFileModel struct {
	models.Model
	GroupID    string    `gorm:"size:64;index" json:"groupId"`
	FileName   string    `gorm:"size:256;index" json:"fileName"`
	UploaderID string    `gorm:"size:64;index" json:"uploaderId"`
	MessageID  string    `gorm:"size:64" json:"messageId"`
	IsPinned   bool      `gorm:"default:false" json:"isPinned"`
	UploadedAt time.Time `json:"uploadedAt"`
}


3、 历史公告
group_announcement_model.go
package group_models

import "beaver/common/models"

// GroupAnnouncementModel 群公告历史表（保留每次公告变更记录）
type GroupAnnouncementModel struct {
	models.Model
	GroupID   string `gorm:"size:64;index" json:"groupId"`
	Content   string `gorm:"type:text" json:"content"`
	CreatedBy string `gorm:"size:64;index" json:"createdBy"`
	Pinned    bool   `gorm:"default:false" json:"pinned"`
	Version   int64  `gorm:"not;default:0;index" json:"version"`
}

// 默认表名：group_announcement_models（如需自定义，可在迁移层统一处理）
