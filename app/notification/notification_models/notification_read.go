package notification_models

import (
	"beaver/common/models"
	"time"
)

// NotificationRead 用户分类已读游标：记录用户对每个分类的最后查看时间
// 用于支持增量同步和状态管理
type NotificationRead struct {
	models.Model
	UserID     string     `gorm:"size:64;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:1;not null" json:"userId"`
	Category   string     `gorm:"size:32;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:2;not null" json:"category"`
	Version    int64      `gorm:"not null;default:0;index" json:"version"` // 游标版本（按用户+分类递增），便于幂等/同步
	LastReadAt *time.Time `json:"lastReadAt,omitempty"`                    // 最后查看该分类的时间
}
