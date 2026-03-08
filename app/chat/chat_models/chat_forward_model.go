package chat_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ForwardContent 转发内容集合类型，实现 GORM 的 Scan 和 Value 接口
type ForwardContent []ChatMessage

func (m ForwardContent) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *ForwardContent) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

// ChatForwardDetail 合并转发详情表（冷热分离，存储大块JSON数据）
type ChatForward struct {
	models.Model
	RecordID string         `gorm:"size:64;uniqueIndex" json:"recordId"` // 聚合ID
	Content  ForwardContent `gorm:"type:json" json:"content"`            // 序列化后的消息数组快照
}
