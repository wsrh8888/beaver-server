package auth

import (
	"context"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"

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
	for _, group := range []string{"desktop", "mobile"} {
		key := fmt.Sprintf("user_authentication_session:%s:%s", req.UserID, group)
		l.svcCtx.Redis.Del(key)
	}

	l.Logger.Infof("用户 %s 登出成功,时间: %s", req.UserID, time.Now().Format("2006-01-02 15:04:05"))
	return &types.LogoutRes{}, nil
}
