package auth

import (
	"context"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		logger: logger.New("logout"),
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.LogoutRes, error) {
	for _, group := range []string{"desktop", "mobile"} {
		key := fmt.Sprintf("user_authentication_session:%s:%s", req.UserID, group)
		l.svcCtx.Redis.Del(key)
	}

	l.logger.Info(model.LogMsg{
		Text: "用户登出成功",
		Data: map[string]interface{}{"userId": req.UserID},
	})
	return &types.LogoutRes{}, nil
}
