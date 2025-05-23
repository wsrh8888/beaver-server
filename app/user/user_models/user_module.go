package user_models

import (
	"beaver/common/models"
)

type UserModel struct {
	models.Model
	UUID        string `gorm:"size:64;unique;index"`
	NickName    string `gorm:"size:32;index"`                                          // 昵称
	Password    string `gorm:"size:128"`                                               // 存储加密后的密码
	Avatar      string `gorm:"size:256;default:e7be4283-dc79-4db7-b65c-aa335b90bcfb" ` // 头像
	Abstract    string `gorm:"size:128"`                                               // 个性签名
	Phone       string `gorm:"size:11;index"`                                          // 手机号
	Status      int8   `gorm:"default:1"`                                              // 1:正常 2:禁用 3:删除
	LastLoginIP string `gorm:"size:39"`                                                // 最后登录IP
	Source      int32  // 注册来源
}
