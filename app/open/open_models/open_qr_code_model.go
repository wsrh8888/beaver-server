package open_models

import (
	"time"

	"gorm.io/gorm"
)

// 扫码登录记录表
type OpenQrCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	SceneID   string         `gorm:"column:scene_id;type:varchar(64);not null;index;comment:场景ID" json:"sceneId"`
	AppID     string         `gorm:"column:app_id;type:varchar(64);not null;index;comment:应用ID" json:"appId"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);default:'';comment:扫码用户ID(扫码后填充)" json:"userId"`
	Status    int            `gorm:"column:status;type:tinyint;not null;default:0;comment:状态:0-等待扫码,1-已扫码,2-已确认,3-已取消,4-已过期" json:"status"`
	ExpiresAt time.Time      `gorm:"column:expires_at;type:datetime;not null;comment:过期时间" json:"expiresAt"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index;comment:删除时间" json:"-"`
}

func (OpenQrCode) TableName() string {
	return "open_qr_code"
}
