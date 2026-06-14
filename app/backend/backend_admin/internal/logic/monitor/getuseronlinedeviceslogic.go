package monitor

import (
	"context"

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOnlineDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOnlineDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOnlineDevicesLogic {
	return &GetUserOnlineDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserOnlineDevicesLogic) GetUserOnlineDevices(req *types.GetUserOnlineDevicesReq) (*types.GetUserOnlineDevicesRes, error) {
	res, err := l.svcCtx.AuthRpc.ListUserDevices(l.ctx, &auth_rpc.ListUserDevicesReq{
		UserId: req.UserID,
	})
	if err != nil {
		l.Errorf("查询用户设备失败 userId=%s: %v", req.UserID, err)
		return nil, err
	}

	list := make([]types.UserOnlineDeviceItem, 0, len(res.List))
	for _, d := range res.List {
		list = append(list, types.UserOnlineDeviceItem{
			DeviceID:        d.DeviceId,
			DeviceType:      d.DeviceType,
			DeviceName:      d.DeviceName,
			DeviceOs:        d.DeviceOs,
			DeviceModel:     d.DeviceModel,
			DeviceOsVersion: d.DeviceOsVersion,
			LastLoginTime:   d.LastLoginTime,
			LastLoginIp:     d.LastLoginIp,
			IsOnline:        d.IsOnline,
		})
	}

	return &types.GetUserOnlineDevicesRes{List: list}, nil
}
