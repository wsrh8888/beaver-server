package open_models

import (
	"gorm.io/gorm"
)

// ==================== Robot 配置表 ====================

// OpenAppRobot 应用智能机器人配置表（对标飞书开放平台）
// 注意：这是应用级的智能机器人配置（AI 对话），不是群内的 Webhook 推送机器人（bot）
type OpenAppRobot struct {
	gorm.Model
	AppID string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`

	// 消息接收配置
	MessageReceiveURL string `gorm:"type:varchar(512);comment:消息接收回调地址"`

	// 功能开关
	EnableSingleChat int `gorm:"type:tinyint;default:1;comment:是否启用单聊 1是 0否"`
	EnableGroupChat  int `gorm:"type:tinyint;default:1;comment:是否启用群聊 1是 0否"`
	EnableAtMention  int `gorm:"type:tinyint;default:1;comment:是否允许@提及 1是 0否"`

	// 自动回复配置
	AutoReplyRules string `gorm:"type:text;comment:自动回复规则(JSON)"`
	WelcomeMessage string `gorm:"type:text;comment:欢迎语"`
	CommandPrefix  string `gorm:"type:varchar(10);default:/;comment:命令前缀"`

	// 事件订阅
	EventSubscriptions string `gorm:"type:text;comment:订阅的事件列表(JSON数组)"`

	// 安全配置
	VerifyToken string `gorm:"type:varchar(128);comment:验证 Token"`
	EncryptKey  string `gorm:"type:varchar(128);comment:加密密钥（可选）"`
	IPWhitelist string `gorm:"type:text;comment:IP白名单(JSON数组)"`
}
