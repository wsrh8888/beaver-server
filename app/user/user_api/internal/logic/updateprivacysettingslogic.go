package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePrivacySettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePrivacySettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePrivacySettingsLogic {
	return &UpdatePrivacySettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePrivacySettingsLogic) UpdatePrivacySettings(req *types.UpdatePrivacySettingsReq) (resp *types.UpdatePrivacySettingsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
