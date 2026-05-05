package app

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Bot 配置
func NewGetBotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotConfigLogic {
	return &GetBotConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBotConfigLogic) GetBotConfig(req *types.GetBotConfigReq) (resp *types.GetBotConfigRes, err error) {
	// todo: add your logic here and delete this line

	return
}
