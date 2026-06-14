package oauth

import (
	"encoding/json"
	"errors"
	"time"

	"beaver/app/open/constants"
	"beaver/app/open/open_models"
	util "beaver/utils/uuid"

	"gorm.io/gorm"
)

const oauthCodeTTL = 5 * time.Minute

const (
	QrStatusWaiting   = 0
	QrStatusScanned   = 1
	QrStatusConfirmed = 2
	QrStatusCancelled = 3
	QrStatusExpired   = 4
)

func QrStatusText(status int) string {
	switch status {
	case QrStatusWaiting:
		return "waiting"
	case QrStatusScanned:
		return "scanned"
	case QrStatusConfirmed:
		return "confirmed"
	case QrStatusCancelled:
		return "cancelled"
	case QrStatusExpired:
		return "expired"
	default:
		return "waiting"
	}
}

// Qrcode 扫码登录状态机（scan / confirm / cancel 共用）
type Qrcode struct {
	db *gorm.DB
}

func NewQrcode(db *gorm.DB) *Qrcode {
	return &Qrcode{db: db}
}

func (q *Qrcode) LoadScene(sceneID string) (*open_models.OpenOAuthQrCode, error) {
	if sceneID == "" {
		return nil, errors.New("sceneId 不能为空")
	}

	var qrCode open_models.OpenOAuthQrCode
	if err := q.db.Where("scene_id = ?", sceneID).First(&qrCode).Error; err != nil {
		return nil, errors.New("二维码不存在或已过期")
	}

	if time.Now().After(qrCode.ExpiresAt) {
		_ = q.db.Model(&qrCode).Update("status", QrStatusExpired).Error
		return nil, errors.New("二维码已过期")
	}

	return &qrCode, nil
}

func (q *Qrcode) MarkScanned(sceneID, userID string) error {
	qrCode, err := q.LoadScene(sceneID)
	if err != nil {
		return err
	}

	if qrCode.Status == QrStatusConfirmed {
		return errors.New("二维码已确认")
	}
	if qrCode.Status == QrStatusCancelled {
		return errors.New("二维码已取消")
	}
	if qrCode.Status == QrStatusScanned && qrCode.UserID != "" && qrCode.UserID != userID {
		return errors.New("二维码已被其他用户扫描")
	}

	return q.db.Model(qrCode).Updates(map[string]interface{}{
		"user_id": userID,
		"status":  QrStatusScanned,
	}).Error
}

func (q *Qrcode) Confirm(sceneID, userID string) error {
	qrCode, err := q.LoadScene(sceneID)
	if err != nil {
		return err
	}

	if qrCode.Status == QrStatusConfirmed {
		return errors.New("二维码已确认")
	}
	if qrCode.Status == QrStatusCancelled {
		return errors.New("二维码已取消")
	}
	if qrCode.Status == QrStatusScanned && qrCode.UserID != "" && qrCode.UserID != userID {
		return errors.New("二维码已被其他用户扫描")
	}

	var app open_models.OpenApp
	if err := q.db.Where("app_id = ? AND status = ?", qrCode.AppID, 1).First(&app).Error; err != nil {
		return errors.New("应用不存在或未启用")
	}

	if err := q.db.Model(qrCode).Updates(map[string]interface{}{
		"user_id": userID,
		"status":  QrStatusConfirmed,
	}).Error; err != nil {
		return errors.New("更新扫码状态失败")
	}

	return q.createOAuthCode(qrCode.AppID, userID, "pc_scan", sceneID)
}

func (q *Qrcode) Cancel(sceneID, userID string) error {
	qrCode, err := q.LoadScene(sceneID)
	if err != nil {
		return err
	}

	if qrCode.Status == QrStatusConfirmed {
		return errors.New("二维码已确认，无法取消")
	}
	if qrCode.Status == QrStatusCancelled {
		return nil
	}
	if qrCode.Status == QrStatusScanned && qrCode.UserID != "" && qrCode.UserID != userID {
		return errors.New("无权取消该扫码会话")
	}

	return q.db.Model(qrCode).Updates(map[string]interface{}{
		"user_id": userID,
		"status":  QrStatusCancelled,
	}).Error
}

func (q *Qrcode) FindConfirmedCode(sceneID string, qrCode *open_models.OpenOAuthQrCode) (string, error) {
	var oauthCode open_models.OpenOAuthCode
	err := q.db.Where(
		"app_id = ? AND user_id = ? AND scene = ? AND used = ? AND state = ?",
		qrCode.AppID, qrCode.UserID, "pc_scan", false, sceneID,
	).Order("id DESC").First(&oauthCode).Error
	if err != nil {
		return "", err
	}
	return oauthCode.Code, nil
}

func (q *Qrcode) createOAuthCode(appID, userID, scene, sceneRef string) error {
	var oauthConfig open_models.OpenAppOAuth
	scope := ""
	if err := q.db.Where("app_id = ?", appID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scope = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scope = string(data)
	}

	code := util.NewV4().String()
	record := open_models.OpenOAuthCode{
		Code:      code,
		AppID:     appID,
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(oauthCodeTTL).Unix(),
		Scene:     scene,
		State:     sceneRef,
	}
	if err := q.db.Create(&record).Error; err != nil {
		return errors.New("生成授权码失败")
	}
	return nil
}
