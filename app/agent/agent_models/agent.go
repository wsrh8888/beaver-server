package agent_models

import (
	"beaver/common/models"
)

// Agent 用户侧的 AI 助手身份（名字、头像用于会话列表展示）。
type Agent struct {
	models.Model
	AgentID string `gorm:"size:64;uniqueIndex;not null;comment:Agent业务ID" json:"agentId"`
	UserID  string `gorm:"size:64;index;not null;comment:创建者用户ID" json:"userId"`
	Name    string `gorm:"size:64;not null;comment:展示名称(会话列表标题)" json:"name"`
	Avatar  string `gorm:"size:500;comment:头像URL(会话列表头像)" json:"avatar"`
	Status  int8   `gorm:"type:tinyint;not null;default:1;index;comment:0停用1启用" json:"status"`
}

func (Agent) TableName() string {
	return "agents"
}
