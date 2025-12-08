package emoji_models

import "beaver/common/models"

type EmojiPackage struct {
	models.Model
	PackageID   string `gorm:"column:package_id;size:64;uniqueIndex" json:"packageId"` // 全局唯一ID
	Title       string `json:"title"`                                                  // 表情包名称
	CoverFile   string `json:"coverFile"`                                              // 表情包封面文件
	UserID      string `json:"userID"`                                                 // 创建者ID（用户ID或官方ID）
	Description string `json:"description"`                                            // 表情包描述
	Type        string `json:"type"`                                                   // 类型：official-官方，user-用户自定义
	Status      int8   `gorm:"default:1" json:"status"`                                // 状态：1=正常 2=审核中 3=违规禁用
	Version     int64  `gorm:"not null;default:0;index" json:"version"`                // 表情包版本号，每次修改递增
}
