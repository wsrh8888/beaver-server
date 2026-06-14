package open_models

import (
	"time"

	"gorm.io/gorm"
)

// OpenOAuthQrCode 扫码登录记录表
type OpenOAuthQrCode struct {
	ID        uint           `gorm:"primarykey"`
	SceneID   string         `gorm:"column:scene_id;type:varchar(64);not null;index"`
	AppID     string         `gorm:"column:app_id;type:varchar(64);not null;index"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);default:''"`
	Status    int            `gorm:"column:status;type:tinyint;not null;default:0;comment:0-等待扫码,1-已扫码,2-已确认,3-已取消,4-已过期"`
	ExpiresAt time.Time      `gorm:"column:expires_at;type:datetime;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
}

func (OpenOAuthQrCode) TableName() string {
	return "open_oauth_qr_codes"
}
