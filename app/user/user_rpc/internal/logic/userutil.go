package logic

import (
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"
)

func toUserInfo(user user_models.UserModel) *user_rpc.UserInfo {
	return &user_rpc.UserInfo{
		UserId:    user.UserID,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Version:   user.Version,
		Email:     user.Email,
		Abstract:  user.Abstract,
		Phone:     user.Phone,
		Status:    int32(user.Status),
		Source:    user.Source,
		UserType:  int32(user.UserType),
		CreatedAt: time.Time(user.CreatedAt).Format(time.RFC3339),
		UpdatedAt: time.Time(user.UpdatedAt).Format(time.RFC3339),
	}
}
