package logic

import (
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.LogoutRes, error) {
	// 删除用户token
	tokenKey := fmt.Sprintf("login_%s", req.UserID)
	if err := l.svcCtx.Redis.Del(tokenKey).Err(); err != nil {
		l.Logger.Errorf("删除用户token失败: %v", err)
		return nil, fmt.Errorf("登出失败")
	}

	// 删除设备信息
	deviceKey := fmt.Sprintf("device_%s", req.UserID)
	if err := l.svcCtx.Redis.Del(deviceKey).Err(); err != nil {
		l.Logger.Errorf("删除设备信息失败: %v", err)
		return nil, fmt.Errorf("登出失败")
	}

	// 记录登出日志
	l.Logger.Infof("用户 %s 登出成功,时间: %s", req.UserID, time.Now().Format("2006-01-02 15:04:05"))

	return &types.LogoutRes{}, nil
}
