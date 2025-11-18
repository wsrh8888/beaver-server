package user_models

import (
	"beaver/common/models"
)

type UserModel struct {
	models.Model
	UUID     string `gorm:"size:64;unique;index"`
	NickName string `gorm:"size:32;index"`                                         // 昵称
	Password string `gorm:"size:128"`                                              // 存储加密后的密码
	Email    string `gorm:"size:128;index"`                                        // 邮箱
	Avatar   string `gorm:"size:256;default:a9de5548bef8c10b92428fff61275c72.png"` // 头像文件ID
	Abstract string `gorm:"size:128"`                                              // 个性签名
	Phone    string `gorm:"size:11;index"`                                         // 手机号
	Status   int8   `gorm:"default:1"`                                             // 1:正常 2:禁用 3:删除
	Gender   int8   `gorm:"default:3"`                                             // 1:男 2:女 3:未知
	Source   int32  // 注册来源
	Version  int64  `gorm:"not null;default:0;index"` // 版本号（用户独立递增）
}
