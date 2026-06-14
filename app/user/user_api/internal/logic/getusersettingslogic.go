package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"

	"gorm.io/gorm"
)

type GetUserSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettingsLogic {
	return &GetUserSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSettingsLogic) GetUserSettings(req *types.GetUserSettingsReq) (*types.GetUserSettingsRes, error) {
	setting, err := l.getOrCreateUserSetting(req.UserID)
	if err != nil {
		return nil, err
	}

	defaults := user_models.DefaultUserSetting(req.UserID).SettingInfo
	info := setting.SettingInfo
	if info == nil {
		info = defaults
	}
	if info.Privacy == nil {
		info.Privacy = defaults.Privacy
	}
	if info.Notification == nil {
		info.Notification = defaults.Notification
	}
	if info.Keyboard == nil {
		info.Keyboard = defaults.Keyboard
	}

	return &types.GetUserSettingsRes{
		Privacy: types.GetUserSettingsPrivacyItem{
			AllowFriendRequest: info.Privacy.AllowFriendRequest,
			ShowOnlineStatus:   info.Privacy.ShowOnlineStatus,
			AllowSearchByPhone: info.Privacy.AllowSearchByPhone,
			AllowSearchByEmail: info.Privacy.AllowSearchByEmail,
		},
		Notification: types.GetUserSettingsNotificationItem{
			NotifyFriendRequest: info.Notification.NotifyFriendRequest,
			NotifyGroupMessage:  info.Notification.NotifyGroupMessage,
			NotifyMoment:        info.Notification.NotifyMoment,
		},
		Keyboard: types.GetUserSettingsKeyboardItem{
			Screenshot:   info.Keyboard.Screenshot,
			ToggleWindow: info.Keyboard.ToggleWindow,
			SendMessage:  info.Keyboard.SendMessage,
		},
	}, nil
}

func (l *GetUserSettingsLogic) getOrCreateUserSetting(userID string) (*user_models.UserSettingModel, error) {
	var setting user_models.UserSettingModel
	err := l.svcCtx.DB.Where("user_id = ?", userID).First(&setting).Error
	if err == nil {
		return &setting, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	setting = user_models.DefaultUserSetting(userID)
	if err := l.svcCtx.DB.Create(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
