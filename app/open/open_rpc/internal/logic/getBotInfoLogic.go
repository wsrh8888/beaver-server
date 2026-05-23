package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotInfoLogic {
	return &GetBotInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetBotInfo 获取推送机器人信息
func (l *GetBotInfoLogic) GetBotInfo(in *open_rpc.GetBotInfoReq) (*open_rpc.GetBotInfoRes, error) {
	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("bot_id = ?", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	return &open_rpc.GetBotInfoRes{
		Id:    uint32(bot.ID),
		BotId: bot.BotID,
		Token: bot.Token,
		Security: &open_rpc.BotSecurity{
			KeywordsEnabled:    bot.Security.KeywordsEnabled,
			Keywords:           bot.Security.Keywords,
			IpWhitelistEnabled: bot.Security.IPWhitelistEnabled,
			IpWhitelist:        bot.Security.IPWhitelist,
			SignatureEnabled:   bot.Security.SignatureEnabled,
			SignatureSecret:    bot.Security.SignatureSecret,
		},
	}, nil
}
