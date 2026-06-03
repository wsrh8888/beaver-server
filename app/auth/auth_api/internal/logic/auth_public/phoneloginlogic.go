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
	"beaver/utils/device"
	"beaver/utils/jwts"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PhoneLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhoneLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneLoginLogic {
	return &PhoneLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneLoginLogic) PhoneLogin(req *types.PhoneLoginReq) (*types.PhoneLoginRes, error) {
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
		return nil, errors.New("密码错误")
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

	now := time.Now()
	credential.LastLoginAt = &now
	credential.LoginCount++
	_ = l.svcCtx.DB.Save(&credential).Error

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

	return &types.PhoneLoginRes{Token: token, UserID: userInfo.UserId}, nil
}
