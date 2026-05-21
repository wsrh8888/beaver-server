package open_models

import "gorm.io/gorm"

type OpenBotModel struct {
	gorm.Model
	AppID             string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`
	BotID             string `gorm:"type:varchar(64);uniqueIndex;comment:Bot的UserID"`
	MessageReceiveURL string `gorm:"type:varchar(512);comment:消息接收回调地址"`
	Name        string `gorm:"type:varchar(100);comment:Bot名称"`
	Avatar      string `gorm:"type:varchar(500);comment:Bot头像URL"`
	Description string `gorm:"type:text;comment:Bot简介"`
	UsageGuide        string `gorm:"type:text;comment:使用说明"`
	EnableSingleChat  int    `gorm:"type:tinyint;default:1;comment:是否启用单聊 1是 0否"`
	EnableGroupChat   int    `gorm:"type:tinyint;default:1;comment:是否启用群聊 1是 0否"`
	EnableAtMention   int    `gorm:"type:tinyint;default:1;comment:是否允许@提及 1是 0否"`
	EnableMenu        int    `gorm:"type:tinyint;default:0;comment:是否启用自定义菜单 1是 0否"`
	MenuItems         string `gorm:"type:text;comment:菜单项配置(JSON)"`
	AutoReplyRules    string `gorm:"type:text;comment:自动回复规则(JSON)"`
	Commands          string `gorm:"type:text;comment:命令列表(JSON)"`
	Status            int    `gorm:"type:tinyint;default:1;comment:状态 1启用 0禁用"`
}
