package update_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
)

// 架构版本配置信息
type StrategyInfo struct {
	ArchitectureID uint `json:"architectureId"` // 架构ID
	VersionID      uint `json:"versionId"`      // 版本ID
	ForceUpdate    bool `json:"forceUpdate"`    // 是否强制更新
	IsActive       bool `json:"isActive"`       // 是否启用
}

type Strategy []StrategyInfo

// Scan 方法 - 从数据库读取数据
func (c *Strategy) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

// Value 方法 - 存储到数据库
func (c *Strategy) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

// 应用城市策略模型
type UpdateStrategy struct {
	models.Model
	AppID    string    `json:"appId" gorm:"size:64;index"` // 应用ID（关联 UpdateApp.AppID）
	CityID   string    `json:"cityId" gorm:"size:32"`      // 城市代码
	Strategy *Strategy `json:"strategy"`                   // 城市策略配置
	IsActive bool      `json:"isActive"`                   // 是否启用
}
