package friend_models

import (
	"beaver/common/models"
	"fmt"

	"gorm.io/gorm"
)

/**
 * @description: 好友表
 */
type FriendModel struct {
	models.Model
	FriendID       string `gorm:"column:friend_id;size:64;uniqueIndex" json:"friendId"`
	FriendshipID   string `gorm:"size:64;unique;index" json:"friendshipId"` // 好友关系唯一ID (min_user_id_max_user_id)
	SendUserID     string `gorm:"size:64;index" json:"sendUserId"`          // 发起验证方的 UserID
	RevUserID      string `gorm:"size:64;index" json:"revUserId"`           // 接收验证方的 UserID
	SendUserNotice string `gorm:"size: 128" json:"sendUserNotice"`          //发起验证方备注
	RevUserNotice  string `gorm:"size: 128" json:"revUserNotice"`           //接收验证方备注
	Source         string `gorm:"size: 32" json:"source"`                   // 好友关系来源：qrcode/search/group/recommend
	IsDeleted      bool   `gorm:"not null;default:false" json:"isDeleted"`  // 标记用户是否删除会话
	Version        int64  `gorm:"not null;default:0;index"`                 // 版本号（用于数据同步）
}

// GenerateFriendshipID 生成好友关系唯一ID (使用min_user_id_max_user_id格式)
func GenerateFriendshipID(userA, userB string) string {
	if userA < userB {
		return fmt.Sprintf("%s_%s", userA, userB)
	}
	return fmt.Sprintf("%s_%s", userB, userA)
}

// BeforeCreate GORM钩子：在创建前生成FriendshipID
func (f *FriendModel) BeforeCreate(tx *gorm.DB) error {
	if f.FriendshipID == "" {
		f.FriendshipID = GenerateFriendshipID(f.SendUserID, f.RevUserID)
	}
	return nil
}

// BeforeUpdate GORM钩子：在更新前确保FriendshipID存在
func (f *FriendModel) BeforeUpdate(tx *gorm.DB) error {
	if f.FriendshipID == "" {
		f.FriendshipID = GenerateFriendshipID(f.SendUserID, f.RevUserID)
	}
	return nil
}

// MigrateFriendshipIDs 迁移现有数据，为空的FriendshipID字段生成值
func MigrateFriendshipIDs(db *gorm.DB) error {
	var friends []FriendModel
	err := db.Where("friendship_id IS NULL OR friendship_id = ''").Find(&friends).Error
	if err != nil {
		return err
	}

	for _, friend := range friends {
		friendshipID := GenerateFriendshipID(friend.SendUserID, friend.RevUserID)
		err = db.Model(&friend).Update("friendship_id", friendshipID).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FriendModel) IsFriend(db *gorm.DB, A, B string) bool {
	err := db.Take(&f, "((send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?) ) and is_deleted = ?", A, B, B, A, false).Error
	return err == nil
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
