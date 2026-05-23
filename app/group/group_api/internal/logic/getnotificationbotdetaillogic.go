package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNotificationBotDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知机器人详情
func NewGetNotificationBotDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotificationBotDetailLogic {
	return &GetNotificationBotDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotificationBotDetailLogic) GetNotificationBotDetail(req *types.GetNotificationBotDetailReq) (resp *types.GetNotificationBotDetailRes, err error) {
	// 查询机器人
	var bot group_models.GroupBotModel
	if err := l.svcCtx.DB.First(&bot, req.ID).Error; err != nil {
		return nil, err
	}

	return &types.GetNotificationBotDetailRes{
		ID:            int64(bot.Id),
		Name:          bot.Name,
		Description:   bot.Description,
		Avatar:        bot.Avatar,
		WebhookURL:    bot.WebhookURL,
		Type:          bot.Type,
		Status:        bot.Status,
		CreatorUserID: bot.CreatorUserID,
		CreatedAt:     time.Time(bot.CreatedAt).Unix(),
	}, nil
}
