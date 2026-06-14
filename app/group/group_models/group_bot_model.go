package group_models

import "beaver/common/models"

// GroupBotModel 群内机器人展示模型（只存群内特有信息）
// 基础信息（昵称、头像等）从 user 表获取
// 安全凭证（Token、签名等）在 open_bots 表管理
type GroupBotModel struct {
	models.Model
	GroupID   string `gorm:"size:128;index;not null" json:"groupId"`     // 群组ID
	BotID     string `gorm:"size:128;uniqueIndex;not null" json:"botId"` // 机器人用户ID（关联 users.user_id）
	Status    int    `gorm:"default:1" json:"status"`                    // 1启用 0禁用
	Type      string `gorm:"size:32;default:'custom'" json:"type"`       // 集成类型：custom/github/gitlab/jenkins/grafana/prometheus
	CreatorID string `gorm:"size:128" json:"creatorId"`                  // 创建者用户ID
}
