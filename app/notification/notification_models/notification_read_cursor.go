package notification_models

import (
	"beaver/common/models"
	"time"
)

// NotificationReadCursor 按类型的已读游标：用于批量标记“读至某事件”
type NotificationReadCursor struct {
	models.Model
	UserID       string     `gorm:"size:64;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:1;not null" json:"userId"`
	Category     string     `gorm:"size:32;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:2;not null" json:"category"`
	Version      int64      `gorm:"not null;default:0;index" json:"version"` // 游标版本（按用户+分类递增），便于幂等/同步
	LastEventID  string     `gorm:"size:64" json:"lastEventId"`
	LastReadAt   *time.Time `json:"lastReadAt,omitempty"`
	LastReadTime int64      `gorm:"index" json:"lastReadTime"` // 冗余时间戳，便于范围过滤
}
