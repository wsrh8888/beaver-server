package bot

import (
	"context"
	"errors"

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
	// TODO: OpenBotModel 已重构为群机器人模型，应用维度的 Bot 配置功能暂时禁用
	return nil, errors.New("Bot 配置功能暂未实现")
}
