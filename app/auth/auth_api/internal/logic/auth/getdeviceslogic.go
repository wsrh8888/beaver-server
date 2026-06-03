package auth

import (
	"context"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"

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
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&devices).Error; err != nil {
		l.Errorf("查询设备失败: %v", err)
		return nil, err
	}

	list := make([]types.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		list = append(list, types.DeviceInfo{
			DeviceID:      d.DeviceID,
			DeviceType:    d.DeviceType,
			DeviceOS:      d.DeviceOS,
			DeviceName:    d.DeviceName,
			LastLoginTime: d.LastLoginTime.Format("2006-01-02 15:04:05"),
			IsActive:      d.IsActive,
			LastLoginIP:   d.LastLoginIP,
		})
	}
	return &types.GetDevicesRes{Devices: list}, nil
}
