package bot

import (
	"context"
	"errors"

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
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// TODO: OpenBotModel 已重构为群机器人模型，应用维度的 Bot 配置功能暂时禁用
	return nil, errors.New("Bot 配置功能暂未实现")
}
