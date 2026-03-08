package doc_models

import (
	"beaver/common/models"
)

// DocVersionModel 文档修订历史/快照保存
// 对标大厂设计：记录每一次重大保存版本，支持一键还原。
type DocVersionModel struct {
	models.Model
	DocID   string `gorm:"size:64;index" json:"docId"` // 文档ID
	FileKey string `gorm:"size:64" json:"fileKey"`     // 该版本对应的文件存储 key
	UserID  string `gorm:"size:64" json:"userId"`      // 提交该版本的用户
	Version int64  `gorm:"default:0" json:"version"`   // 快照对应的逻辑版本号
	Remark  string `gorm:"size:255" json:"remark"`     // 备注说明 (如：手动快照、自动保存)
}
