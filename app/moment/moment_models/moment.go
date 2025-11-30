package moment_models

import (
	"beaver/common/models"
	"beaver/common/models/ctype"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FileInfo 结构体，用于存储文件的信息
type FileInfo struct {
	FileKey string        `json:"fileKey"` // 文件名
	Type    ctype.MsgType `json:"type"`    // 文件类型：使用MsgType枚举，与消息系统保持一致
}

type Files []FileInfo

// Value converts the Files slice to a JSON-encoded string for database storage
func (f *Files) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan converts a JSON-encoded string from the database to a Files slice
func (f *Files) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

/**
 * @description: 动态表
 */
type MomentModel struct {
	models.Model
	UUID      string `gorm:"size:64;uniqueIndex;not null" json:"uuid"`      // 全局唯一ID (UUID，跨库同步用)
	UserID    string `gorm:"size:64;not null;index" json:"userId"`          // 用户Id (索引，提升查询性能)
	Content   string `gorm:"type:text;not null" json:"content"`             // 动态内容
	Files     *Files `gorm:"type:longtext" json:"files"`                    // 文件信息（JSON数组）
	IsDeleted bool   `gorm:"not null;default:false;index" json:"isDeleted"` // 软删除标记 (索引)
}
