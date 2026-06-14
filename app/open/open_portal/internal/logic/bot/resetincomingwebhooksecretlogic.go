package bot

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetIncomingWebhookSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetIncomingWebhookSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetIncomingWebhookSecretLogic {
	return &ResetIncomingWebhookSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetIncomingWebhookSecretLogic) ResetIncomingWebhookSecret(req *types.ResetIncomingWebhookSecretReq) (resp *types.ResetIncomingWebhookSecretRes, err error) {
	if req.ID == "" {
		return nil, errors.New("id 不能为空")
	}

	botID, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil || botID == 0 {
		return nil, errors.New("id 无效")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", botID).First(&bot).Error; err != nil {
		return nil, errors.New("记录不存在")
	}
	if bot.AppID == "" {
		return nil, errors.New("无法操作非 Portal 创建的 Bot")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", bot.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	secretRes, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{Id: uint32(bot.ID)})
	if err != nil {
		return nil, errors.New("重置密钥失败")
	}

	return &types.ResetIncomingWebhookSecretRes{
		Secret: secretRes.SignatureSecret,
	}, nil
}
