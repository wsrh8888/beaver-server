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
	MomentID   string `gorm:"column:moment_id;size:64;uniqueIndex;not null" json:"momentId"` // 全局唯一ID
	UserID     string `gorm:"size:64;not null;index" json:"userId"`                          // 用户Id
	Content    string `gorm:"type:text;not null" json:"content"`                             // 动态内容
	Files      *Files `gorm:"type:longtext" json:"files"`                                    // 文件信息（JSON数组）
	Visibility int8   `gorm:"not null;default:0" json:"visibility"`                          // 可见性：0=所有人 1=仅好友 2=仅自己
	AllowList  string `gorm:"type:text" json:"allowList"`                                    // 白名单用户ID，逗号分隔（visibility=3时生效）
	BlockList  string `gorm:"type:text" json:"blockList"`                                    // 黑名单用户ID，逗号分隔（不让谁看）
	IsDeleted  bool   `gorm:"not null;default:false;index" json:"isDeleted"`                 // 软删除标记
}
