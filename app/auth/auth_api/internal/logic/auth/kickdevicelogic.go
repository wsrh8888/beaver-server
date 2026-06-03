package auth

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickDeviceLogic {
	return &KickDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickDeviceLogic) KickDevice(req *types.KickDeviceReq) (*types.KickDeviceRes, error) {
	var device auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", req.UserID, req.DeviceID).
		First(&device).Error; err != nil {
		return nil, errors.New("设备不存在")
	}

	redisKey := fmt.Sprintf("user_authentication_session:%s:%s", req.UserID, device.DeviceType)
	l.svcCtx.Redis.Del(redisKey)

	l.svcCtx.DB.Model(&device).Update("is_active", false)

	if _, err := l.svcCtx.NotificationRpc.DisablePushToken(l.ctx, &notification_rpc.DisablePushTokenReq{
		UserId:   req.UserID,
		DeviceId: req.DeviceID,
	}); err != nil {
		l.Errorf("禁用 Push Token 失败: userId=%s, deviceId=%s, err=%v", req.UserID, req.DeviceID, err)
	}

	go func() {
		payload := map[string]interface{}{
			"command":  wsCommandConst.USER_PROFILE,
			"type":     wsTypeConst.UserKickReceive,
			"senderId": req.UserID,
			"targetId": req.UserID,
			"body": map[string]interface{}{
				"deviceId": req.DeviceID,
				"reason":   "kicked_by_user",
			},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
	}()

	return &types.KickDeviceRes{}, nil
}
