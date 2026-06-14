package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CheckDeviceActiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckDeviceActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckDeviceActiveLogic {
	return &CheckDeviceActiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckDeviceActiveLogic) CheckDeviceActive(in *auth_rpc.CheckDeviceActiveReq) (*auth_rpc.CheckDeviceActiveRes, error) {
	if in.UserId == "" || in.DeviceId == "" {
		return nil, errors.New("userId 和 deviceId 不能为空")
	}

	var device auth_models.AuthDeviceModel
	err := l.svcCtx.DB.Where("user_id = ? AND device_id = ? AND is_active = ?",
		in.UserId, in.DeviceId, true).First(&device).Error
	if err == gorm.ErrRecordNotFound {
		return &auth_rpc.CheckDeviceActiveRes{Active: false}, nil
	}
	if err != nil {
		l.Errorf("校验设备状态失败: userId=%s, deviceId=%s, err=%v", in.UserId, in.DeviceId, err)
		return nil, errors.New("校验设备状态失败")
	}

	return &auth_rpc.CheckDeviceActiveRes{Active: true}, nil
}
