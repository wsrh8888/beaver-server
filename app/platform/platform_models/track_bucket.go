package platform_models

import "beaver/common/models"

// TrackBucket 埋点与日志的 Bucket 注册表
type TrackBucket struct {
	models.Model
	Name        string `json:"name" gorm:"size:64"`
	Description string `json:"description" gorm:"type:text"`
	BucketID    string `json:"bucketId" gorm:"column:bucket_id;uniqueIndex;size:64"`
	Kind        string `json:"kind" gorm:"index;size:16"` // track=埋点 log=日志
	CreateUser  string `json:"createUser" gorm:"size:64"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
}
