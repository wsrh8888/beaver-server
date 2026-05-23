package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取机器人详情
func NewGetBotDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotDetailLogic {
	return &GetBotDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBotDetailLogic) GetBotDetail(req *types.GetBotDetailReq) (resp *types.GetBotDetailRes, err error) {
	// 1. 查询群内机器人信息（通过 bot_id 查询）
	var bot group_models.GroupBotModel
	if err := l.svcCtx.DB.Where("bot_id = ?", req.BotID).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	// 2. 通过 user_rpc 获取用户基础信息（昵称、头像）
	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: bot.BotID,
	})
	if err != nil || userRes.UserInfo == nil {
		return nil, err
	}

	// 3. 通过 open_rpc 获取 Webhook Token 和安全设置
	botInfoRes, err := l.svcCtx.OpenRpc.GetBotInfo(l.ctx, &open_rpc.GetBotInfoReq{
		BotId: bot.BotID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 拼接完整 Webhook URL
	fullWebhookURL := fmt.Sprintf("%s/api/webhook/%s?token=%s", l.svcCtx.Config.Domain, botInfoRes.Token, botInfoRes.Token)

	return &types.GetBotDetailRes{
		BotID:         bot.BotID,
		Name:          userRes.UserInfo.NickName,
		Description:   userRes.UserInfo.Abstract, // 用户简介
		Avatar:        userRes.UserInfo.Avatar,
		WebhookURL:    fullWebhookURL,
		Type:          bot.Type,
		Status:        bot.Status,
		CreatorUserID: bot.CreatorID,
		CreatedAt:     time.Time(bot.CreatedAt).Unix(),
		Security: types.BotSecurity{
			KeywordsEnabled:    botInfoRes.Security.KeywordsEnabled,
			Keywords:           botInfoRes.Security.Keywords,
			IPWhitelistEnabled: botInfoRes.Security.IpWhitelistEnabled,
			IPWhitelist:        botInfoRes.Security.IpWhitelist,
			SignatureEnabled:   botInfoRes.Security.SignatureEnabled,
			SignatureSecret:    botInfoRes.Security.SignatureSecret,
		},
	}, nil
}
