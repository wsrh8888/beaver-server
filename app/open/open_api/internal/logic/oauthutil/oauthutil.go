package oauthutil

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"beaver/app/open/constants"
	"beaver/app/open/open_models"
	util "beaver/utils/uuid"

	"gorm.io/gorm"
)

const oauthCodeTTL = 5 * time.Minute

func VerifyApp(db *gorm.DB, appID, appSecret string) (*open_models.OpenApp, error) {
	if appID == "" || appSecret == "" {
		return nil, errors.New("应用凭证不完整")
	}

	var app open_models.OpenApp
	err := db.Where("app_id = ? AND app_secret = ? AND status = ?", appID, appSecret, 1).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}

	return &app, nil
}

// VerifyAppForCodeExchange 授权码换 token：appSecret 可选（PKCE 公开客户端可不传）
func VerifyAppForCodeExchange(db *gorm.DB, appID, appSecret string) (*open_models.OpenApp, error) {
	if appID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	query := db.Where("app_id = ? AND status = ?", appID, 1)
	if appSecret != "" {
		query = query.Where("app_secret = ?", appSecret)
	}
	if err := query.First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}
	return &app, nil
}

func ParseScopes(scopeStr string) []string {
	scopeStr = strings.TrimSpace(scopeStr)
	if scopeStr == "" {
		return nil
	}
	if strings.HasPrefix(scopeStr, "[") {
		var scopes []string
		if err := json.Unmarshal([]byte(scopeStr), &scopes); err == nil {
			return scopes
		}
	}
	return strings.FieldsFunc(scopeStr, func(r rune) bool {
		return r == ' ' || r == ','
	})
}

func HasScope(granted []string, required constants.ScopeType) bool {
	target := string(required)
	for _, s := range granted {
		if s == target {
			return true
		}
	}
	return false
}

func RequireScope(record *open_models.OpenOAuthToken, required constants.ScopeType) error {
	if !HasScope(ParseScopes(record.Scope), required) {
		return errors.New("权限不足: 缺少 " + string(required))
	}
	return nil
}

func RequireUserToken(record *open_models.OpenOAuthToken) error {
	if record.UserID == "" {
		return errors.New("需要用户授权令牌")
	}
	return nil
}

func ValidateAccessTokenWithScopes(db *gorm.DB, authorization string, required ...constants.ScopeType) (*open_models.OpenOAuthToken, error) {
	record, err := ValidateAccessToken(db, authorization)
	if err != nil {
		return nil, err
	}
	for _, scope := range required {
		if err := RequireScope(record, scope); err != nil {
			return nil, err
		}
	}
	return record, nil
}

func resolveAppScope(db *gorm.DB, appID string) string {
	return ResolveAppScope(db, appID)
}

// ResolveAppScope 获取应用授权 scope（JSON 数组字符串）
func ResolveAppScope(db *gorm.DB, appID string) string {
	var oauthConfig open_models.OpenAppOAuth
	if err := db.Where("app_id = ?", appID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		return oauthConfig.SupportedScopes
	}

	scopes := []string{
		string(constants.ScopeUserProfileRead),
		string(constants.ScopeUserAvatarRead),
	}
	data, _ := json.Marshal(scopes)
	return string(data)
}

func CreateOAuthCode(db *gorm.DB, appID, userID, scene string) (code string, expireIn int64, err error) {
	if appID == "" || userID == "" {
		return "", 0, errors.New("参数不完整")
	}

	var app open_models.OpenApp
	if err := db.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error; err != nil {
		return "", 0, errors.New("应用不存在或未启用")
	}

	code = util.NewV4().String()
	expiresAt := time.Now().Add(oauthCodeTTL).Unix()
	record := open_models.OpenOAuthCode{
		Code:      code,
		AppID:     appID,
		UserID:    userID,
		Scope:     resolveAppScope(db, appID),
		ExpiresAt: expiresAt,
		Scene:     scene,
	}

	if err := db.Create(&record).Error; err != nil {
		return "", 0, errors.New("生成授权码失败")
	}

	return code, int64(oauthCodeTTL.Seconds()), nil
}

func FindOAuthCode(db *gorm.DB, appID, code string) (*open_models.OpenOAuthCode, error) {
	var record open_models.OpenOAuthCode
	err := db.Where("code = ? AND app_id = ?", code, appID).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("授权码无效")
		}
		return nil, errors.New("查询授权码失败")
	}
	return &record, nil
}

func ValidateOAuthCode(record *open_models.OpenOAuthCode) error {
	if record.Used {
		return errors.New("授权码已使用")
	}
	if time.Now().Unix() > record.ExpiresAt {
		return errors.New("授权码已过期")
	}
	return nil
}

func MarkOAuthCodeUsed(db *gorm.DB, record *open_models.OpenOAuthCode) error {
	return db.Model(record).Update("used", true).Error
}

func RevokeOAuthToken(db *gorm.DB, token string) error {
	if token == "" {
		return errors.New("token 不能为空")
	}

	result := db.Where("token = ? OR refresh_token = ?", token, token).Delete(&open_models.OpenOAuthToken{})
	if result.Error != nil {
		return errors.New("撤销令牌失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("令牌不存在")
	}
	return nil
}

func ParseBearerToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return authorization
}

func ValidateAccessToken(db *gorm.DB, authorization string) (*open_models.OpenOAuthToken, error) {
	token := ParseBearerToken(authorization)
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
