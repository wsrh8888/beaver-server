package emoji_models

import (
	"beaver/common/models"
)

// 表情
type Emoji struct {
	models.Model
	UUID    string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	FileKey string `json:"fileKey"`                                 // 文件Key
	Title   string `json:"title"`                                   // 表情名称
	Status  int8   `gorm:"default:1" json:"status"`                 // 状态：1=正常 2=审核中 3=违规禁用
	Version int64  `gorm:"not null;default:0;index" json:"version"` //基于表递增
}
