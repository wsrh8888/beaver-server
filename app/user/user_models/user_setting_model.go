package user_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PrivacySetting 隐私设置
type PrivacySetting struct {
	AllowFriendRequest bool `json:"allowFriendRequest"`
	ShowOnlineStatus   bool `json:"showOnlineStatus"`
	AllowSearchByPhone bool `json:"allowSearchByPhone"`
	AllowSearchByEmail bool `json:"allowSearchByEmail"`
}

// NotificationSetting 通知策略
type NotificationSetting struct {
	NotifyFriendRequest bool `json:"notifyFriendRequest"`
	NotifyGroupMessage  bool `json:"notifyGroupMessage"`
	NotifyMoment        bool `json:"notifyMoment"`
}

// SettingInfo 用户设置详情（一人一条，JSON 内包含各模块）
type SettingInfo struct {
	Privacy      *PrivacySetting      `json:"privacy"`
	Notification *NotificationSetting `json:"notification"`
}

func (s *SettingInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SettingInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// UserSettingModel 用户设置，一个用户一行
type UserSettingModel struct {
	models.Model
	UserID      string       `gorm:"size:64;uniqueIndex" json:"userId"`
	SettingInfo *SettingInfo `gorm:"type:longtext" json:"settingInfo"`
}

func (UserSettingModel) TableName() string {
	return "user_settings"
}

func DefaultUserSetting(userID string) UserSettingModel {
	return UserSettingModel{
		UserID: userID,
		SettingInfo: &SettingInfo{
			Privacy: &PrivacySetting{
				AllowFriendRequest: true,
				ShowOnlineStatus:   true,
				AllowSearchByPhone: true,
				AllowSearchByEmail: true,
			},
			Notification: &NotificationSetting{
				NotifyFriendRequest: true,
				NotifyGroupMessage:  true,
				NotifyMoment:        true,
			},
		},
	}
}
