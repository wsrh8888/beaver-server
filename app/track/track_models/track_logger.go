package track_models

import (
	"beaver/common/models"

	"gorm.io/datatypes"
)

// 日志埋点事件表 - 用于系统监控和问题排查
type TrackLogger struct {
	models.Model
	Level       string         `json:"level" gorm:"index;size:16"`                             // 日志级别
	Data        datatypes.JSON `json:"data"`                                                   // 日志数据
	BucketID    string         `json:"bucketId" gorm:"index"`                                  // Bucket Id
	BucketModel *TrackBucket   `gorm:"foreignKey:BucketID;references:UUID" json:"bucketModel"` // Bucket关联信息
	Timestamp   int64          `json:"timestamp" gorm:"index"`                                 // 日志时间戳
}
