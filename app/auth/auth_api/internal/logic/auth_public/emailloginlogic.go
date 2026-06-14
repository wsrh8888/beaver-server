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
	"beaver/utils/authlock"
	"beaver/utils/device"
	"beaver/utils/jwts"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type EmailLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		ctx:    ctx,
		logger: logger.New("email_login"),
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

	_ = device.UpsertOnLogin(l.svcCtx.DB, userInfo.UserId, req.DeviceID, ctxUAProfile(l.ctx), ctxClientIP(l.ctx))

	l.logger.Info(model.LogMsg{
		Text: "邮箱验证码登录成功",
		Data: map[string]interface{}{
			"userId":      userInfo.UserId,
			"deviceGroup": deviceGroup,
		},
	})

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
	return authlock.VerifyStoredCode(l.ctx, l.svcCtx.Redis, codeKey, codeType, email, code)
}
