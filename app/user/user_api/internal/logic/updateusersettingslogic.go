package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserSettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserSettingsLogic {
	return &UpdateUserSettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserSettingsLogic) UpdateUserSettings(req *types.UpdateUserSettingsReq) (resp *types.UpdateUserSettingsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
