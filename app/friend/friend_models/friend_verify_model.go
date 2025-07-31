package friend_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
)

/**
 * @description: 好友验证
 */
type FriendVerifyModel struct {
	models.Model
	RevUserModel  user_models.UserModel `gorm:"foreignkey:RevUserID;references:UUID" json:"-"`
	SendUserModel user_models.UserModel `gorm:"foreignkey:SendUserID;references:UUID" json:"-"`
	SendUserID    string                `gorm:"size:64;index" json:"sendUserId"` // 使用 VARCHAR(64)
	RevUserID     string                `gorm:"size:64;index" json:"revUserId"`  // 使用 VARCHAR(64)
	SendStatus    int8                  `json:"sendStatus"`                      // 发起方状态 0:未处理 1:已通过 2:已拒绝 3: 忽略 4:删除
	RevStatus     int8                  `json:"revStatus"`                       // 接收方状态 0:未处理 1:已通过 2:已拒绝 3: 忽略 4:删除
	Message       string                `gorm:"size: 128" json:"message"`        // 附加消息
	Source        string                `gorm:"size: 32" json:"source"`          // 添加好友来源：qrcode/search/group/recommend
}
