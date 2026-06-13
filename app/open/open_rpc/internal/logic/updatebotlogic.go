package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotLogic {
	return &UpdateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBotLogic) UpdateBot(in *open_rpc.UpdateBotReq) (*open_rpc.UpdateBotRes, error) {
	if in.BotId == "" {
		return nil, errors.New("botId 不能为空")
	}
	if in.Security == nil {
		return &open_rpc.UpdateBotRes{}, nil
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("bot_id = ? AND status = 1", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	security := open_models.OpenBotSecurity{
		KeywordsEnabled:    in.Security.KeywordsEnabled,
		IPWhitelistEnabled: in.Security.IpWhitelistEnabled,
		SignatureEnabled:   in.Security.SignatureEnabled,
	}
	if in.Security.KeywordsEnabled {
		security.Keywords = in.Security.Keywords
	}
	if in.Security.IpWhitelistEnabled {
		security.IPWhitelist = in.Security.IpWhitelist
	}
	if in.Security.SignatureSecret != "" {
		security.SignatureSecret = in.Security.SignatureSecret
	} else {
		security.SignatureSecret = bot.Security.SignatureSecret
	}

	if err := l.svcCtx.DB.Model(&bot).Update("security", security).Error; err != nil {
		return nil, errors.New("更新机器人安全设置失败")
	}

	return &open_rpc.UpdateBotRes{}, nil
}
