package track_models

import "beaver/common/models"

// TrackBucket 表示不同的bucket/命名空间
type TrackBucket struct {
	models.Model
	Name        string `json:"name" gorm:"size:64"`                                  // 软件名称，如 "即时通讯系统", "客户关系管理系统"
	Description string `json:"description" gorm:"type:text"`                         // 软件描述
	BucketID    string `json:"bucketId" gorm:"column:bucket_id;uniqueIndex;size:64"` // 唯一标识符
	CreateUser  string `json:"CreateUser" gorm:"size:64"`                            // 负责人
	IsActive    bool   `json:"isActive" gorm:"default:true"`                         // 是否激活
}
