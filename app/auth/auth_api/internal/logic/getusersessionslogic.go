package logic

import (
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSessionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSessionsLogic {
	return &GetUserSessionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSessionsLogic) GetUserSessions(req *types.GetUserSessionsReq) (*types.GetUserSessionsRes, error) {
	// 从Redis获取设备信息
	deviceKey := fmt.Sprintf("device_%s", req.UserID)
	deviceInfo, err := l.svcCtx.Redis.Get(deviceKey).Result()
	if err != nil {
		l.Logger.Errorf("获取设备信息失败: %v", err)
		return &types.GetUserSessionsRes{
			Sessions: []types.SessionInfo{},
		}, nil
	}

	// 解析设备信息
	var session types.SessionInfo
	if err := json.Unmarshal([]byte(deviceInfo), &session); err != nil {
		l.Logger.Errorf("解析设备信息失败: %v", err)
		return &types.GetUserSessionsRes{
			Sessions: []types.SessionInfo{},
		}, nil
	}

	// 格式化最后活跃时间
	lastActive, err := time.Parse(time.RFC3339, session.LastActive)
	if err == nil {
		session.LastActive = lastActive.Format("2006-01-02 15:04:05")
	}

	return &types.GetUserSessionsRes{
		Sessions: []types.SessionInfo{session},
	}, nil
}
