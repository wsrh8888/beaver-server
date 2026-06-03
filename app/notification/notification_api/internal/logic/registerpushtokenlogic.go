package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterPushTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterPushTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterPushTokenLogic {
	return &RegisterPushTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterPushTokenLogic) RegisterPushToken(req *types.RegisterPushTokenReq) (*types.RegisterPushTokenRes, error) {
	if req.DeviceID == "" {
		return nil, errors.New("deviceId 不能为空")
	}
	if req.PushToken == "" {
		return nil, errors.New("pushToken 不能为空")
	}

	platform := strings.ToLower(strings.TrimSpace(req.PushPlatform))
	if platform != "fcm" && platform != "apns" {
		return nil, errors.New("pushPlatform 仅支持 fcm 或 apns")
	}

	res, err := l.svcCtx.AuthRpc.CheckDeviceActive(l.ctx, &auth_rpc.CheckDeviceActiveReq{
		UserId:   req.UserID,
		DeviceId: req.DeviceID,
	})
	if err != nil {
		l.Errorf("校验设备状态失败: userId=%s, deviceId=%s, err=%v", req.UserID, req.DeviceID, err)
		return nil, errors.New("校验设备状态失败")
	}
	if res == nil || !res.Active {
		return nil, errors.New("设备未登录或已失效，请先登录")
	}

	var row notification_models.PushRegistrationModel
	err = l.svcCtx.DB.Where("user_id = ? AND device_id = ?", req.UserID, req.DeviceID).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		err = l.svcCtx.DB.Create(&notification_models.PushRegistrationModel{
			UserID:       req.UserID,
			DeviceID:     req.DeviceID,
			PushToken:    req.PushToken,
			PushPlatform: platform,
			Enabled:      true,
		}).Error
	} else if err == nil {
		err = l.svcCtx.DB.Model(&row).Updates(map[string]interface{}{
			"push_token":    req.PushToken,
			"push_platform": platform,
			"enabled":       true,
			"updated_at":    time.Now(),
		}).Error
	}
	if err != nil {
		l.Errorf("注册 Push Token 失败: userId=%s, deviceId=%s, err=%v", req.UserID, req.DeviceID, err)
		return nil, errors.New("注册 Push Token 失败")
	}

	l.Infof("注册 Push Token: userId=%s, deviceId=%s, platform=%s", req.UserID, req.DeviceID, platform)
	return &types.RegisterPushTokenRes{}, nil
}
