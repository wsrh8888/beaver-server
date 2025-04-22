package logic

import (
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type TerminateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTerminateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TerminateSessionLogic {
	return &TerminateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TerminateSessionLogic) TerminateSession(req *types.TerminateSessionReq) (*types.TerminateSessionRes, error) {
	// 删除设备信息
	deviceKey := fmt.Sprintf("device_%s", req.UserID)
	if err := l.svcCtx.Redis.Del(deviceKey).Err(); err != nil {
		l.Logger.Errorf("删除设备信息失败: %v", err)
		return nil, fmt.Errorf("终止会话失败")
	}

	// 删除token
	tokenKey := fmt.Sprintf("login_%s", req.UserID)
	if err := l.svcCtx.Redis.Del(tokenKey).Err(); err != nil {
		l.Logger.Errorf("删除token失败: %v", err)
		return nil, fmt.Errorf("终止会话失败")
	}

	// 记录终止日志
	l.Logger.Infof("用户 %s 终止会话成功,设备ID: %s,时间: %s",
		req.UserID,
		req.DeviceID,
		time.Now().Format("2006-01-02 15:04:05"))

	return &types.TerminateSessionRes{}, nil
}
