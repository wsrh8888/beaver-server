package user_models

import (
	"beaver/common/models"
)

// 用户类型常量
const (
	UserTypeNormal int8 = 1 // 普通用户
	UserTypeBot    int8 = 2 // 推送机器人（通知机器人）
	UserTypeRobot  int8 = 3 // 智能机器人（AI 对话机器人）
)

// UserModel 用户基础信息模型（包括普通用户和机器人）
type UserModel struct {
	models.Model
	UserID   string `gorm:"size:64;uniqueIndex" json:"userId"`
	UserType int8   `gorm:"default:1;index" json:"userType"`                       // 用户类型：1普通用户 2bot 3robot
	NickName string `gorm:"size:32;index" json:"nickName"`                         // 昵称
	Avatar   string `gorm:"size:256;default:a9de5548bef8c10b92428fff61275c72.png"` // 头像文件ID
	Abstract string `gorm:"size:128" json:"abstract"`                              // 个性签名
	Email    string `gorm:"size:128;index" json:"email"`                           // 邮箱（普通用户有，机器人为NULL）
	Phone    string `gorm:"size:11;index" json:"phone"`                            // 手机号（普通用户有，机器人为NULL）
	Status   int8   `gorm:"default:1" json:"status"`                               // 1:正常 2:禁用 3:删除
	Gender   int8   `gorm:"default:3" json:"gender"`                               // 1:男 2:女 3:未知（仅普通用户）
	Source   int32  `json:"source"`                                                // 注册来源
	Version  int64  `gorm:"not null;default:0;index" json:"version"`               // 版本号（用户独立递增）
}
