package group_models

import "beaver/common/models"

// GroupBotModel 群内机器人展示模型（只存展示信息）
type GroupBotModel struct {
	models.Model
	GroupID       string `gorm:"size:128;index;not null" json:"groupId"`
	BotID         uint   `gorm:"uniqueIndex;not null" json:"botId"` // 关联 open_bots.id
	Name          string `gorm:"size:100" json:"name"`
	Description   string `gorm:"size:500" json:"description"`
	Avatar        string `gorm:"size:256" json:"avatar"`
	WebhookURL    string `gorm:"size:512" json:"webhookUrl"`
	Status        int    `gorm:"default:1" json:"status"`
	Type          string `gorm:"size:32;default:'custom'" json:"type"`
	CreatorUserID string `gorm:"size:128" json:"creatorUserId"`
}
