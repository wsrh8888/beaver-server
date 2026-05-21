package group_models

import "beaver/common/models"

type GroupNotificationBotModel struct {
	models.Model
	GroupID       string `gorm:"size:128;index;not null" json:"groupId"`
	WebhookID     uint   `gorm:"uniqueIndex;not null" json:"webhookId"` // open_bots.id
	BotUserID     string `gorm:"size:128;not null" json:"botUserId"`
	Name          string `gorm:"size:100" json:"name"`
	Description   string `gorm:"size:500" json:"description"`
	Avatar        string `gorm:"size:256" json:"avatar"`
	WebhookURL    string `gorm:"size:512" json:"webhookUrl"`
	Status        int    `gorm:"default:1" json:"status"`
	Type          string `gorm:"size:32;default:'custom'" json:"type"`
	CreatorUserID string `gorm:"size:128" json:"creatorUserId"`
}
