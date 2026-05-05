package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新 Bot 配置
func NewUpdateBotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotConfigLogic {
	return &UpdateBotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBotConfigLogic) UpdateBotConfig(req *types.UpdateBotConfigReq) (resp *types.UpdateBotConfigRes, err error) {
	// todo: add your logic here and delete this line

	return
}
