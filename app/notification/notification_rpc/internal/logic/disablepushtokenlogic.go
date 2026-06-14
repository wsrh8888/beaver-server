package logic

import (
	"context"
	"errors"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisablePushTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisablePushTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisablePushTokenLogic {
	return &DisablePushTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DisablePushTokenLogic) DisablePushToken(in *notification_rpc.DisablePushTokenReq) (*notification_rpc.DisablePushTokenRes, error) {
	if in.UserId == "" || in.DeviceId == "" {
		return nil, errors.New("userId 和 deviceId 不能为空")
	}

	result := l.svcCtx.DB.Model(&notification_models.PushRegistrationModel{}).
		Where("user_id = ? AND device_id = ?", in.UserId, in.DeviceId).
		Update("enabled", false)
	if result.Error != nil {
		l.Errorf("禁用 Push Token 失败: userId=%s, deviceId=%s, err=%v", in.UserId, in.DeviceId, result.Error)
		return nil, errors.New("禁用 Push Token 失败")
	}

	return &notification_rpc.DisablePushTokenRes{Success: true}, nil
}
