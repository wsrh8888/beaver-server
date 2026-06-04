package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	uuidUtil "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetBotSecretLogic) ResetBotSecret(in *open_rpc.ResetBotSecretReq) (*open_rpc.ResetBotSecretRes, error) {
	if in.Id == 0 {
		return nil, errors.New("id 不能为空")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	security := bot.Security
	security.SignatureEnabled = true
	security.SignatureSecret = uuidUtil.NewV4().String()

	if err := l.svcCtx.DB.Model(&bot).Update("security", security).Error; err != nil {
		return nil, errors.New("重置密钥失败")
	}

	return &open_rpc.ResetBotSecretRes{
		SignatureSecret: security.SignatureSecret,
	}, nil
}
