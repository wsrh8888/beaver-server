package backend_models

import (
	"beaver/common/models"
)

/**
 * @description: 管理员用户表
 * 注意：使用软删除，Status 只保留 1:正常 2:禁用，删除使用 DeletedAt
 */
type AdminUser struct {
	models.Model
	UserID      string `gorm:"column:user_id;size:64;uniqueIndex;comment:管理员ID"`                 // 管理员ID（唯一标识）
	NickName    string `gorm:"size:32;index;comment:昵称"`                                         // 昵称
	Password    string `gorm:"size:128;comment:加密后的密码"`                                          // 存储加密后的密码
	Avatar      string `gorm:"size:256;default:a9de5548bef8c10b92428fff61275c72.png;comment:头像"` // 头像文件ID
	Abstract    string `gorm:"size:128;comment:个性签名"`                                            // 个性签名
	Phone       string `gorm:"size:11;unique;index;comment:手机号"`                                 // 手机号（唯一）
	Status      int8   `gorm:"default:1;index;comment:状态"`                                       // 1:正常 2:禁用（删除使用软删除）
	LastLoginAt int64  `gorm:"index;comment:最后登录时间"`                                             // 最后登录时间戳
	CreatedBy   string `gorm:"size:64;index;comment:创建者管理员ID"`                                   // 创建者管理员ID
}
