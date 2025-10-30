package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户设备列表
func NewGetDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDevicesLogic {
	return &GetDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDevicesLogic) GetDevices(req *types.GetDevicesReq) (resp *types.GetDevicesRes, err error) {
	// 查询用户的所有设备
	var devices []user_models.UserDeviceModel
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&devices).Error
	if err != nil {
		logx.Errorf("查询用户设备失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var deviceInfos []types.DeviceInfo
	for _, device := range devices {
		deviceInfos = append(deviceInfos, types.DeviceInfo{
			DeviceID:      device.DeviceID,
			DeviceType:    device.DeviceType,
			DeviceOS:      device.DeviceOS,
			DeviceName:    device.DeviceName,
			LastLoginTime: device.LastLoginTime.Format("2006-01-02 15:04:05"),
			IsActive:      device.IsActive,
			LastLoginIP:   device.LastLoginIP,
		})
	}

	return &types.GetDevicesRes{
		Devices: deviceInfos,
	}, nil
}
