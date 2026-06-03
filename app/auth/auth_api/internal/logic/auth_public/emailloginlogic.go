package auth_public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/middleware/ua"
	"beaver/utils/device"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type EmailLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailLoginLogic) EmailLogin(req *types.EmailLoginReq) (*types.EmailLoginRes, error) {
	if err := l.verifyEmailCode(req.Email, req.Code, "login"); err != nil {
		return nil, err
	}

	searchRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Email,
		Type:    "email",
	})
	var userInfo *user_rpc.UserInfo
	if err != nil {
		userInfo, err = l.createUserByEmail(req.Email)
		if err != nil {
			return nil, errors.New("用户注册失败")
		}
	} else {
		userInfo = searchRes.UserInfo
	}

	preciseType, _ := l.ctx.Value(ua.KeyDeviceType).(string)
	deviceGroup, _ := l.ctx.Value(ua.KeyDeviceGroup).(string)
	if preciseType == "" || preciseType == device.DeviceUnknown {
		return nil, errors.New("不支持的设备类型")
	}
	if req.DeviceID == "" {
		return nil, errors.New("无法识别的物理设备，请联系管理员")
	}

	key := fmt.Sprintf("user_authentication_session:%s:%s", userInfo.UserId, deviceGroup)
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: userInfo.NickName,
		UserID:   userInfo.UserId,
		DeviceID: req.DeviceID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.New("服务内部异常")
	}

	loginInfo, _ := json.Marshal(map[string]interface{}{
		"token": token, "device_id": req.DeviceID, "device_type": preciseType,
		"device_group": deviceGroup, "login_time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err = l.svcCtx.Redis.Set(key, string(loginInfo), time.Hour*time.Duration(l.svcCtx.Config.Auth.AccessExpire)).Err(); err != nil {
		return nil, errors.New("服务内部异常")
	}

	var credential auth_models.AuthCredentialModel
	if l.svcCtx.DB.Take(&credential, "user_id = ?", userInfo.UserId).Error == nil {
		now := time.Now()
		credential.LastLoginAt = &now
		credential.LoginCount++
		_ = l.svcCtx.DB.Save(&credential).Error
	}

	now := time.Now()
	var dev auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", userInfo.UserId, req.DeviceID).First(&dev).Error; err == gorm.ErrRecordNotFound {
		_ = l.svcCtx.DB.Create(&auth_models.AuthDeviceModel{
			UserID: userInfo.UserId, DeviceID: req.DeviceID, DeviceType: deviceGroup, DeviceOS: preciseType,
			LastLoginTime: now, IsActive: true,
		}).Error
	} else if err == nil {
		_ = l.svcCtx.DB.Model(&dev).Updates(map[string]interface{}{
			"device_type": deviceGroup, "device_os": preciseType,
			"last_login_time": now, "is_active": true, "updated_at": now,
		}).Error
	}

	return &types.EmailLoginRes{Token: token, UserID: userInfo.UserId}, nil
}

func (l *EmailLoginLogic) createUserByEmail(email string) (*user_rpc.UserInfo, error) {
	nickName := fmt.Sprintf("用户%s", email[:strings.Index(email, "@")])
	createRes, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateReq{
		Email: email, NickName: nickName, Source: 2,
	})
	if err != nil {
		return nil, err
	}
	return &user_rpc.UserInfo{UserId: createRes.UserID, NickName: nickName, Email: email}, nil
}

func (l *EmailLoginLogic) verifyEmailCode(email, code, codeType string) error {
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}
	l.svcCtx.Redis.Del(codeKey)
	return nil
}
