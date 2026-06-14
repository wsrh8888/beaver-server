package auth_public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/middleware/ua"
	"beaver/utils/authlock"
	"beaver/utils/device"
	"beaver/utils/jwts"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"beaver/utils/pwd"
)


type PhoneLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewPhoneLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneLoginLogic {
	return &PhoneLoginLogic{
		ctx:    ctx,
		logger: logger.New("phone_login"),
		svcCtx: svcCtx,
	}
}

func (l *PhoneLoginLogic) PhoneLogin(req *types.PhoneLoginReq) (*types.PhoneLoginRes, error) {
	failKey := authlock.LoginFailKey(req.Phone)
	lockKey := authlock.LoginLockKey(req.Phone)
	if err := authlock.CheckLocked(l.ctx, l.svcCtx.Redis, lockKey); err != nil {
		return nil, err
	}

	searchRes, err := l.svcCtx.UserRpc.SearchUser(l.ctx, &user_rpc.SearchUserReq{
		Keyword: req.Phone,
		Type:    "phone",
	})
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	userInfo := searchRes.UserInfo

	var credential auth_models.AuthCredentialModel
	if err := l.svcCtx.DB.Take(&credential, "user_id = ?", userInfo.UserId).Error; err != nil {
		return nil, errors.New("用户凭证不存在")
	}
	if !pwd.CheckPad(credential.Password, req.Password) {
		if lockErr := authlock.RecordFailure(l.ctx, l.svcCtx.Redis, failKey, lockKey, "login", "phone"); lockErr != nil {
			return nil, lockErr
		}
		return nil, errors.New("密码错误")
	}
	authlock.ClearFailures(l.svcCtx.Redis, failKey, lockKey)

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

	now := time.Now()
	credential.LastLoginAt = &now
	credential.LoginCount++
	_ = l.svcCtx.DB.Save(&credential).Error

	_ = device.UpsertOnLogin(l.svcCtx.DB, userInfo.UserId, req.DeviceID, ctxUAProfile(l.ctx), ctxClientIP(l.ctx))

	l.logger.Info(model.LogMsg{
		Text: "手机密码登录成功",
		Data: map[string]interface{}{
			"userId":      userInfo.UserId,
			"deviceGroup": deviceGroup,
		},
	})

	return &types.PhoneLoginRes{Token: token, UserID: userInfo.UserId}, nil
}
