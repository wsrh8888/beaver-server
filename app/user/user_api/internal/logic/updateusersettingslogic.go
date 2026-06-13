package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"gorm.io/gorm"
)

type UpdateUserSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserSettingsLogic {
	return &UpdateUserSettingsLogic{
		ctx:    ctx,
		logger: logger.New("update_user_settings"),
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserSettingsLogic) UpdateUserSettings(req *types.UpdateUserSettingsReq) (*types.UpdateUserSettingsRes, error) {
	setting, err := l.getOrCreateUserSetting(req.UserID)
	if err != nil {
		return nil, err
	}

	if setting.SettingInfo == nil {
		setting.SettingInfo = user_models.DefaultUserSetting(req.UserID).SettingInfo
	}
	if setting.SettingInfo.Privacy == nil {
		setting.SettingInfo.Privacy = user_models.DefaultUserSetting(req.UserID).SettingInfo.Privacy
	}
	if setting.SettingInfo.Notification == nil {
		setting.SettingInfo.Notification = user_models.DefaultUserSetting(req.UserID).SettingInfo.Notification
	}

	if req.Privacy != nil {
		if req.Privacy.AllowFriendRequest != nil {
			setting.SettingInfo.Privacy.AllowFriendRequest = *req.Privacy.AllowFriendRequest
		}
		if req.Privacy.ShowOnlineStatus != nil {
			setting.SettingInfo.Privacy.ShowOnlineStatus = *req.Privacy.ShowOnlineStatus
		}
		if req.Privacy.AllowSearchByPhone != nil {
			setting.SettingInfo.Privacy.AllowSearchByPhone = *req.Privacy.AllowSearchByPhone
		}
		if req.Privacy.AllowSearchByEmail != nil {
			setting.SettingInfo.Privacy.AllowSearchByEmail = *req.Privacy.AllowSearchByEmail
		}
	}

	if req.Notification != nil {
		if req.Notification.NotifyFriendRequest != nil {
			setting.SettingInfo.Notification.NotifyFriendRequest = *req.Notification.NotifyFriendRequest
		}
		if req.Notification.NotifyGroupMessage != nil {
			setting.SettingInfo.Notification.NotifyGroupMessage = *req.Notification.NotifyGroupMessage
		}
		if req.Notification.NotifyMoment != nil {
			setting.SettingInfo.Notification.NotifyMoment = *req.Notification.NotifyMoment
		}
	}

	if err := l.svcCtx.DB.Save(setting).Error; err != nil {
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "用户设置更新成功",
		Data: map[string]interface{}{"userId": req.UserID},
	})
	return &types.UpdateUserSettingsRes{}, nil
}

func (l *UpdateUserSettingsLogic) getOrCreateUserSetting(userID string) (*user_models.UserSettingModel, error) {
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
