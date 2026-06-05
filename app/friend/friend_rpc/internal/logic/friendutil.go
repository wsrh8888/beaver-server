package logic

import (
	"fmt"
	"strconv"
	"time"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"gorm.io/gorm"
)

func friendKey(f friend_models.FriendModel) string {
	if f.FriendID != "" {
		return f.FriendID
	}
	return fmt.Sprintf("%d", f.Id)
}

func toFriendItem(f friend_models.FriendModel) *friend_rpc.FriendItem {
	return &friend_rpc.FriendItem{
		FriendId:       friendKey(f),
		SendUserId:     f.SendUserID,
		RevUserId:      f.RevUserID,
		SendUserNotice: f.SendUserNotice,
		RevUserNotice:  f.RevUserNotice,
		Source:         f.Source,
		IsDeleted:      f.IsDeleted,
		CreatedAt:      time.Time(f.CreatedAt).Format(time.RFC3339),
		UpdatedAt:      time.Time(f.UpdatedAt).Format(time.RFC3339),
	}
}

func toFriendVerifyItem(v friend_models.FriendVerifyModel) *friend_rpc.FriendVerifyItem {
	key := v.VerifyID
	if key == "" {
		key = fmt.Sprintf("%d", v.Id)
	}
	return &friend_rpc.FriendVerifyItem{
		VerifyId:   key,
		SendUserId: v.SendUserID,
		RevUserId:  v.RevUserID,
		SendStatus: int32(v.SendStatus),
		RevStatus:  int32(v.RevStatus),
		Message:    v.Message,
		Source:     v.Source,
		CreatedAt:  time.Time(v.CreatedAt).Format(time.RFC3339),
		UpdatedAt:  time.Time(v.UpdatedAt).Format(time.RFC3339),
	}
}

func findFriend(db *gorm.DB, id string) (*friend_models.FriendModel, error) {
	var f friend_models.FriendModel
	if err := db.Where("friend_id = ?", id).First(&f).Error; err == nil {
		return &f, nil
	}
	if n, err := strconv.ParseUint(id, 10, 64); err == nil {
		if err := db.Where("id = ?", n).First(&f).Error; err == nil {
			return &f, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func blockKey(b friend_models.FriendBlockModel) string {
	if b.BlockID != "" {
		return b.BlockID
	}
	return fmt.Sprintf("%d", b.Id)
}

func toFriendBlockItem(b friend_models.FriendBlockModel) *friend_rpc.FriendBlockItem {
	return &friend_rpc.FriendBlockItem{
		BlockId:       blockKey(b),
		UserId:        b.UserID,
		BlockedUserId: b.BlockedUserID,
		CreatedAt:     time.Time(b.CreatedAt).Format(time.RFC3339),
	}
}

func findFriendBlock(db *gorm.DB, id string) (*friend_models.FriendBlockModel, error) {
	var b friend_models.FriendBlockModel
	if err := db.Where("block_id = ?", id).First(&b).Error; err == nil {
		return &b, nil
	}
	if n, err := strconv.ParseUint(id, 10, 64); err == nil {
		if err := db.Where("id = ?", n).First(&b).Error; err == nil {
			return &b, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func findFriendVerify(db *gorm.DB, id string) (*friend_models.FriendVerifyModel, error) {
	var v friend_models.FriendVerifyModel
	if err := db.Where("verify_id = ?", id).First(&v).Error; err == nil {
		return &v, nil
	}
	if n, err := strconv.ParseUint(id, 10, 64); err == nil {
		if err := db.Where("id = ?", n).First(&v).Error; err == nil {
			return &v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
