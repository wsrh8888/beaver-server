package platform_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
)

type StrategyInfo struct {
	ArchitectureID uint `json:"architectureId"`
	VersionID      uint `json:"versionId"`
	ForceUpdate    bool `json:"forceUpdate"`
	IsActive       bool `json:"isActive"`
}

type Strategy []StrategyInfo

func (c *Strategy) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

func (c *Strategy) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

type UpdateStrategy struct {
	models.Model
	AppID    string    `json:"appId" gorm:"size:64;index"`
	CityID   string    `json:"cityId" gorm:"size:32"`
	Strategy *Strategy `json:"strategy"`
	IsActive bool      `json:"isActive"`
}
