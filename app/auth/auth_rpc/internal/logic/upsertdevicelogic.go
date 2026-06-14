package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpsertDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpsertDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertDeviceLogic {
	return &UpsertDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpsertDeviceLogic) UpsertDevice(in *auth_rpc.UpsertDeviceReq) (*auth_rpc.UpsertDeviceRes, error) {
	if in.UserId == "" || in.DeviceId == "" {
		return nil, errors.New("userId 和 deviceId 不能为空")
	}

	now := time.Now()
	var dev auth_models.AuthDeviceModel
	err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", in.UserId, in.DeviceId).First(&dev).Error
	if err == gorm.ErrRecordNotFound {
		err = l.svcCtx.DB.Create(&auth_models.AuthDeviceModel{
			UserID:        in.UserId,
			DeviceID:      in.DeviceId,
			DeviceType:    in.DeviceType,
			DeviceOS:      in.DeviceOs,
			DeviceName:    in.DeviceName,
			LastLoginTime: now,
			IsActive:      true,
			LastLoginIP:   in.LastLoginIp,
		}).Error
	} else if err == nil {
		updates := map[string]interface{}{
			"device_type":     in.DeviceType,
			"device_os":       in.DeviceOs,
			"device_name":     in.DeviceName,
			"last_login_time": now,
			"is_active":       true,
			"updated_at":      now,
		}
		if in.LastLoginIp != "" {
			updates["last_login_ip"] = in.LastLoginIp
		}
		err = l.svcCtx.DB.Model(&dev).Updates(updates).Error
	}
	if err != nil {
		l.Errorf("登记设备失败: userId=%s, deviceId=%s, err=%v", in.UserId, in.DeviceId, err)
		return nil, errors.New("登记设备失败")
	}
	return &auth_rpc.UpsertDeviceRes{Success: true}, nil
}
