package doc_models

import (
	"beaver/common/models"
)

// DocModel 文档元数据主表 (对标大厂工业级设计：元数据与权限解耦)
type DocModel struct {
	models.Model
	DocID   string `gorm:"size:64;uniqueIndex" json:"docId"` // 业务唯一ID
	FileKey string `gorm:"size:64;index" json:"fileKey"`     // 关联文件系统存储 key
	Title   string `gorm:"size:255;index" json:"title"`      // 文档标题
	Type    int8   `gorm:"index" json:"type"`                // 1:Word, 2:Excel, 3:PPT, 4:Folder
	OwnerID string `gorm:"size:64;index" json:"ownerId"`     // 创建者/所有者

	// 共享策略：1:私有, 2:空间内共享, 3:公开链接
	ShareType int8 `gorm:"default:1" json:"shareType"`

	IsLocked bool  `gorm:"default:false" json:"isLocked"`     // 编辑锁
	Version  int64 `gorm:"not null;default:0" json:"version"` // 乐观锁控制 (处理并发编辑)
	Status   int8  `gorm:"default:1" json:"status"`           // 1:正常, 2:回收站
}
