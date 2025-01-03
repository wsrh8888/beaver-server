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
	RevUserModel   user_models.UserModel `gorm:"foreignkey:RevUserId;references:UserId" json:"-"`
	SendUserModel  user_models.UserModel `gorm:"foreignkey:SendUserId;references:UserId" json:"-"`
	SendUserId     string                `gorm:"size:64;index" json:"sendUserId"`         // 发起验证方的 UserId
	RevUserId      string                `gorm:"size:64;index" json:"revUserId"`          // 接收验证方的 UserId
	SendUserNotice string                `gorm:"size: 128" json:"sendUserNotice"`         //发起验证方备注
	RevUserNotice  string                `gorm:"size: 128" json:"revUserNotice"`          //接收验证方备注
	IsDeleted      bool                  `gorm:"not null;default:false" json:"isDeleted"` // 标记用户是否删除会话
}

/**
 * @description: A -> B SendUserId(A的Id) RevUserId(B的Id) SendUserNotice(A对B的备注） RevUserNotice(B对A的备注)
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
func (f *FriendModel) GetUserNotice(userId string) string {

	if userId == f.SendUserId {
		// 如果我是发起方
		return f.SendUserNotice
	}
	if userId == f.RevUserId {
		// 如果我是接收方
		return f.RevUserNotice
	}
	return ""
}
