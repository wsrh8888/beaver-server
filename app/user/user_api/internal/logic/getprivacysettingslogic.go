package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPrivacySettingsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPrivacySettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivacySettingsLogic {
	return &GetPrivacySettingsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPrivacySettingsLogic) GetPrivacySettings(req *types.GetPrivacySettingsReq) (resp *types.GetPrivacySettingsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
