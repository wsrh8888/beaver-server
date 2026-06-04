package utils

import (
	"errors"
	"strings"
	"time"

	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

func parseBearerToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return authorization
}

func ValidateAppAccessToken(db *gorm.DB, authorization string) (*open_models.OpenOAuthToken, error) {
	token := parseBearerToken(authorization)
	if token == "" {
		return nil, errors.New("缺少访问令牌")
	}

	var record open_models.OpenOAuthToken
	if err := db.Where("token = ?", token).First(&record).Error; err != nil {
		return nil, errors.New("访问令牌无效")
	}
	if time.Now().Unix() > record.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}
	return &record, nil
}

func RequireAppCapability(app *open_models.OpenApp, needRobot, needWebhook bool) error {
	if app.Status != 1 {
		return errors.New("应用未发布或已禁用")
	}
	if needRobot && app.EnableRobot != 1 {
		return errors.New("应用未启用智能机器人能力")
	}
	if needWebhook && app.EnableWebhook != 1 {
		return errors.New("应用未启用 Webhook 能力")
	}
	return nil
}

func LoadAppByID(db *gorm.DB, appID string) (*open_models.OpenApp, error) {
	var app open_models.OpenApp
	if err := db.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在")
	}
	return &app, nil
}
