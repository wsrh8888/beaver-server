package track_models

import (
	"beaver/common/models"

	"gorm.io/datatypes"
)

// 统计埋点事件表 - 关键字段固定，便于分析
type TrackEvent struct {
	models.Model
	EventName   string         `json:"eventName" gorm:"index;size:128"`                            // 事件名称，如 "user_register", "button_click", "page_view"
	Action      string         `json:"action" gorm:"index;size:32"`                                // 操作，如 "click", "view"
	UserID      *string        `json:"userId" gorm:"index;size:64"`                                // 用户ID(可选)
	BucketID    string         `json:"bucketId" gorm:"index"`                                      // Bucket Id
	BucketModel *TrackBucket   `gorm:"foreignKey:BucketID;references:BucketID" json:"bucketModel"` // Bucket关联信息
	Platform    string         `json:"platform" gorm:"size:32"`                                    // 平台，如"ios"、"android"、"web"
	DeviceID    string         `json:"deviceId" gorm:"index;size:64"`                              // 设备ID，用于追踪未登录用户
	Data        datatypes.JSON `json:"data"`                                                       // 事件数据，JSON格式存储所有相关信息
	Timestamp   int64          `json:"timestamp" gorm:"index"`                                     // 日志时间戳

}
