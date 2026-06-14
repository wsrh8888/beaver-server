package open_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// OpenBotSecurity 通知机器人安全设置（JSON 结构）
type OpenBotSecurity struct {
	KeywordsEnabled    bool     `json:"keywordsEnabled"`    // 是否启用关键词校验
	Keywords           []string `json:"keywords"`           // 关键词列表（最多10个）
	IPWhitelistEnabled bool     `json:"ipWhitelistEnabled"` // 是否启用IP白名单
	IPWhitelist        []string `json:"ipWhitelist"`        // IP地址列表
	SignatureEnabled   bool     `json:"signatureEnabled"`   // 是否启用签名校验
	SignatureSecret    string   `json:"signatureSecret"`    // 签名密钥
}

// Value 实现 driver.Valuer 接口
func (s OpenBotSecurity) Value() (driver.Value, error) {
	if !s.KeywordsEnabled && !s.IPWhitelistEnabled && !s.SignatureEnabled {
		return "{}", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *OpenBotSecurity) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// OpenBotModel 推送机器人模型（群内创建的通知机器人，用于接收 Webhook 推送）
// 例如：Jenkins、GitLab、监控告警等第三方服务推送消息到群
type OpenBotModel struct {
	gorm.Model
	AppID   string `gorm:"type:varchar(64);index;comment:开放平台应用ID（Portal 创建时写入）"`
	Name    string `gorm:"type:varchar(100);comment:显示名称"`
	BotID   string `gorm:"type:varchar(64);uniqueIndex;comment:Bot的UserID"`
	GroupID string `gorm:"type:varchar(64);index;not null;comment:目标群组ID"`
	Token   string `gorm:"type:varchar(128);uniqueIndex;comment:Webhook Token（URL参数）"`
	Status  int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`

	// 安全设置（JSON 格式，用于 Webhook 推送时校验）
	Security OpenBotSecurity `gorm:"type:json;comment:安全设置"`
}
