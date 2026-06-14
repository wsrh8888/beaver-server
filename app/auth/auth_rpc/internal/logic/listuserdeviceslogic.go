package logic

import (
	"context"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/core/coreonline"
	"beaver/utils/device"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserDevicesLogic {
	return &ListUserDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserDevicesLogic) ListUserDevices(in *auth_rpc.ListUserDevicesReq) (*auth_rpc.ListUserDevicesRes, error) {
	var devices []auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND is_active = ?", in.UserId, true).
		Order("last_login_time DESC").
		Find(&devices).Error; err != nil {
		l.Errorf("查询用户设备失败 userId=%s: %v", in.UserId, err)
		return nil, err
	}

	onlineDeviceIDs := make(map[string]bool)
	for _, slot := range []string{"desktop", "mobile"} {
		if !coreonline.IsSlotOnline(l.svcCtx.Redis, in.UserId, slot) {
			continue
		}
		deviceID, err := device.SessionDeviceID(l.svcCtx.Redis, in.UserId, slot)
		if err == redis.Nil {
			continue
		}
		if err != nil {
			l.Errorf("读取会话设备失败 userId=%s slot=%s: %v", in.UserId, slot, err)
			return nil, err
		}
		onlineDeviceIDs[deviceID] = true
	}

	list := make([]*auth_rpc.UserDeviceItem, 0, len(devices))
	for _, d := range devices {
		list = append(list, &auth_rpc.UserDeviceItem{
			DeviceId:        d.DeviceID,
			DeviceType:      d.DeviceType,
			DeviceOs:        d.DeviceOS,
			DeviceModel:     d.DeviceModel,
			DeviceOsVersion: d.DeviceOsVersion,
			DeviceName:      d.DeviceName,
			LastLoginTime:   d.LastLoginTime.Format("2006-01-02 15:04:05"),
			LastLoginIp:     d.LastLoginIP,
			IsOnline:        onlineDeviceIDs[d.DeviceID],
		})
	}

	return &auth_rpc.ListUserDevicesRes{List: list}, nil
}
