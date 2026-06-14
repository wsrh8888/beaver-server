package open_models

import (
	"gorm.io/gorm"
)

// OpenAppRobot 应用智能机器人配置（一个 App 对应一个 Robot IM 用户）
type OpenAppRobot struct {
	gorm.Model
	AppID       string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`
	RobotID string `gorm:"column:robot_user_id;type:varchar(64);uniqueIndex;comment:Robot IM 用户ID"`
	RobotName   string `gorm:"type:varchar(100);comment:Robot 昵称"`
	Avatar      string `gorm:"type:varchar(500);comment:Robot 头像"`
	Status      int    `gorm:"type:tinyint;default:1;comment:1启用 0禁用"`

	EnableSingleChat int `gorm:"type:tinyint;default:1;comment:是否启用单聊"`
	EnableGroupChat  int `gorm:"type:tinyint;default:1;comment:是否启用群聊"`
	EnableAtMention  int `gorm:"type:tinyint;default:1;comment:是否允许@提及"`

	WelcomeMessage string `gorm:"type:text;comment:欢迎语"`
	CommandPrefix  string `gorm:"type:varchar(10);default:/;comment:命令前缀"`
}
