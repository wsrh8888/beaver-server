package utils

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type UserCreator interface {
	UserCreate(ctx context.Context, in *user_rpc.UserCreateReq, opts ...grpc.CallOption) (*user_rpc.UserCreateRes, error)
}

func EnsureAppRobot(ctx context.Context, db *gorm.DB, userRpc UserCreator, app *open_models.OpenApp) (*open_models.OpenAppRobot, error) {
	var robot open_models.OpenAppRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		return &robot, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nickName := app.Name
	if nickName == "" {
		nickName = "Robot"
	}

	createRes, err := userRpc.UserCreate(ctx, &user_rpc.UserCreateReq{
		NickName: nickName,
		UserType: int32(user_models.UserTypeRobot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Robot IM 用户失败: %w", err)
	}

	robot = open_models.OpenAppRobot{
		AppID:            app.AppID,
		RobotID:          createRes.UserID,
		RobotName:        nickName,
		Avatar:           app.Icon,
		Status:           1,
		EnableSingleChat: 1,
		EnableGroupChat:  1,
		EnableAtMention:  1,
	}
	if err := db.Save(&robot).Error; err != nil {
		return nil, err
	}
	return &robot, nil
}
