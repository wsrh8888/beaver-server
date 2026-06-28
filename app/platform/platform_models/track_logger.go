package platform_models

import (
	"beaver/common/models"

	"gorm.io/datatypes"
)

// TrackLogger 原始日志表
type TrackLogger struct {
	models.Model
	Level       string         `json:"level" gorm:"index;size:16"`                                 // 日志级别
	Data        datatypes.JSON `json:"data"`                                                       // 日志数据
	BucketID    string         `json:"bucketId" gorm:"index"` // Bucket Id
	Timestamp   int64          `json:"timestamp" gorm:"index"` // 日志时间戳
}
