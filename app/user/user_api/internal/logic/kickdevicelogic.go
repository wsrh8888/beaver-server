package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickDeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 强制下线指定设备
func NewKickDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickDeviceLogic {
	return &KickDeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickDeviceLogic) KickDevice(req *types.KickDeviceReq) (resp *types.KickDeviceRes, err error) {
	// 查询设备，确认归属于当前用户
	var device user_models.UserDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", req.UserID, req.DeviceID).
		First(&device).Error; err != nil {
		return nil, errors.New("设备不存在")
	}

	// 删除该设备类型对应的 Redis 登录态，使 token 立即失效
	redisKey := fmt.Sprintf("login_%s_%s", req.UserID, device.DeviceType)
	l.svcCtx.Redis.Del(redisKey)

	// 标记设备为非活跃
	l.svcCtx.DB.Model(&device).Update("is_active", false)

	// 异步通过 WS 推送强制下线通知（客户端收到后比对 deviceId，执行本地登出）
	go func() {
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.USER_PROFILE, wsTypeConst.UserKickReceive,
			req.UserID, req.UserID, map[string]interface{}{
				"deviceId": req.DeviceID,
				"reason":   "kicked_by_user",
			}, "")
	}()

	return &types.KickDeviceRes{}, nil
}
