package emoji_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// EmojiInfo 表情信息
type EmojiInfo struct {
	Width  int `json:"width"`  // 表情图片宽度
	Height int `json:"height"` // 表情图片高度
}

// Value converts the EmojiInfo to a JSON-encoded string for database storage
func (e *EmojiInfo) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Scan converts a JSON-encoded string from the database to a EmojiInfo
func (e *EmojiInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, e)
}

// 表情
type Emoji struct {
	models.Model
	EmojiID   string     `gorm:"column:emoji_id;size:64;uniqueIndex" json:"emojiId"` // 全局唯一ID
	FileKey   string     `json:"fileKey"`                                            // 文件Key
	Title     string     `json:"title"`                                              // 表情名称
	EmojiInfo *EmojiInfo `gorm:"type:longtext" json:"emojiInfo"`                     // 表情详细信息（JSON格式）
	Status    int8       `gorm:"default:1" json:"status"`                            // 状态：1=正常 2=审核中 3=违规禁用
	Version   int64      `gorm:"not null;default:0;index" json:"version"`            //基于表递增
}
