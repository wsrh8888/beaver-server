package emoji_models

import "beaver/common/models"

type EmojiPackage struct {
	models.Model
	UUID        string `gorm:"size:64;unique;index" json:"uuid"`        // 全局唯一标识符，用于前端同步
	Title       string `json:"title"`                                   // 表情包名称
	CoverFile   string `json:"coverFile"`                               // 表情包封面文件
	UserID      string `json:"userID"`                                  // 创建者ID（用户ID或官方ID）
	Description string `json:"description"`                             // 表情包描述
	Type        string `json:"type"`                                    // 类型：official-官方，user-用户自定义
	Status      int8   `gorm:"default:1" json:"status"`                 // 状态：1=正常 2=审核中 3=违规禁用
	Version     int64  `gorm:"not null;default:0;index" json:"version"` // 表情包UUID版本号， 每次修改内存数据时递增
}
