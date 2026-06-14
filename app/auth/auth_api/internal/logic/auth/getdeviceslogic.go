package auth

import (
	"context"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/core/coreonline"
	"beaver/utils/device"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDevicesLogic {
	return &GetDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDevicesLogic) GetDevices(req *types.GetDevicesReq) (*types.GetDevicesRes, error) {
	var devices []auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).Order("last_login_time DESC").Find(&devices).Error; err != nil {
		l.Errorf("查询设备失败: %v", err)
		return nil, err
	}

	onlineDeviceIDs := make(map[string]bool)
	for _, slot := range []string{"desktop", "mobile"} {
		if !coreonline.IsSlotOnline(l.svcCtx.Redis, req.UserID, slot) {
			continue
		}
		deviceID, err := device.SessionDeviceID(l.svcCtx.Redis, req.UserID, slot)
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		onlineDeviceIDs[deviceID] = true
	}

	list := make([]types.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		if !d.IsActive {
			continue
		}
		list = append(list, types.DeviceInfo{
			DeviceID:        d.DeviceID,
			DeviceType:      d.DeviceType,
			DeviceOS:        d.DeviceOS,
			DeviceModel:     d.DeviceModel,
			DeviceOsVersion: d.DeviceOsVersion,
			DeviceName:      d.DeviceName,
			LastLoginTime:   d.LastLoginTime.Format("2006-01-02 15:04:05"),
			IsOnline:        onlineDeviceIDs[d.DeviceID],
			LastLoginIP:     d.LastLoginIP,
		})
	}
	return &types.GetDevicesRes{Devices: list}, nil
}
