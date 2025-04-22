package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserSettingsLogic {
	return &GetUserSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserSettingsLogic) GetUserSettings(req *types.GetUserSettingsReq) (resp *types.GetUserSettingsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
