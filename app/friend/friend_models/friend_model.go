package friend_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"

	"gorm.io/gorm"
)

/**
 * @description: 好友表
 */
type FriendModel struct {
	models.Model
	RevUserModel   user_models.UserModel `gorm:"foreignkey:RevUserID;references:UUID" json:"-"`
	SendUserModel  user_models.UserModel `gorm:"foreignkey:SendUserID;references:UUID" json:"-"`
	SendUserID     string                `gorm:"size:64;index" json:"sendUserId"`         // 发起验证方的 UserID
	RevUserID      string                `gorm:"size:64;index" json:"revUserId"`          // 接收验证方的 UserID
	SendUserNotice string                `gorm:"size: 128" json:"sendUserNotice"`         //发起验证方备注
	RevUserNotice  string                `gorm:"size: 128" json:"revUserNotice"`          //接收验证方备注
	Source         string                `gorm:"size: 32" json:"source"`                  // 好友关系来源：qrcode/search/group/recommend
	IsDeleted      bool                  `gorm:"not null;default:false" json:"isDeleted"` // 标记用户是否删除会话
}

/**
 * @description: A -> B SendUserID(A的ID) RevUserID(B的ID) SendUserNotice(A对B的备注） RevUserNotice(B对A的备注)
 * @description: B -> A
 */

func (f *FriendModel) IsFriend(db *gorm.DB, A, B string) bool {
	err := db.Take(&f, "((send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?) ) and is_deleted = ?", A, B, B, A, false).Error
	if err != nil {
		return false
	}
	return true
}

/**
 * @description: 获取用户的备注
 * @param {uint} userId
 */
func (f *FriendModel) GetUserNotice(userID string) string {

	if userID == f.SendUserID {
		// 如果我是发起方
		return f.SendUserNotice
	}
	if userID == f.RevUserID {
		// 如果我是接收方
		return f.RevUserNotice
	}
	return ""
}
