package device

import (
	"encoding/json"
	"fmt"
	"time"

	"beaver/app/auth/auth_models"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// UpsertOnLogin 登录成功后登记/更新设备档案
func UpsertOnLogin(db *gorm.DB, userID, deviceID string, profile UAProfile, clientIP string) error {
	now := time.Now()
	var dev auth_models.AuthDeviceModel
	err := db.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&dev).Error
	if err == gorm.ErrRecordNotFound {
		return db.Create(&auth_models.AuthDeviceModel{
			UserID:          userID,
			DeviceID:        deviceID,
			DeviceType:      profile.DeviceGroup,
			DeviceOS:        profile.PreciseType,
			DeviceModel:     profile.Model,
			DeviceOsVersion: profile.OsVersion,
			DeviceName:      profile.DisplayName,
			LastLoginTime:   now,
			IsActive:        true,
			LastLoginIP:     clientIP,
		}).Error
	}
	if err != nil {
		return err
	}

	return db.Model(&dev).Updates(map[string]interface{}{
		"device_type":       profile.DeviceGroup,
		"device_os":         profile.PreciseType,
		"device_model":      profile.Model,
		"device_os_version": profile.OsVersion,
		"device_name":       profile.DisplayName,
		"last_login_time":   now,
		"is_active":         true,
		"last_login_ip":     clientIP,
		"updated_at":        now,
	}).Error
}

// Deactivate 登出或踢下线后标记设备会话失效
func Deactivate(db *gorm.DB, userID, deviceID string) error {
	return db.Model(&auth_models.AuthDeviceModel{}).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Update("is_active", false).Error
}

// SessionDeviceID 读取某槽位当前登录的 deviceId
func SessionDeviceID(rdb *redis.Client, userID, slot string) (string, error) {
	key := fmt.Sprintf("user_authentication_session:%s:%s", userID, slot)
	raw, err := rdb.Get(key).Result()
	if err != nil {
		return "", err
	}
	var info struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return "", err
	}
	return info.DeviceID, nil
}
