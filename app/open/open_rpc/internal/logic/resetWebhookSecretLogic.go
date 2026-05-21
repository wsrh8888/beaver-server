package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetWebhookSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetWebhookSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetWebhookSecretLogic {
	return &ResetWebhookSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetWebhookSecretLogic) ResetWebhookSecret(in *open_rpc.ResetWebhookSecretReq) (*open_rpc.ResetWebhookSecretRes, error) {
	var record open_models.OpenGroupBotModel
	if err := l.svcCtx.DB.First(&record, in.Id).Error; err != nil {
		return nil, errors.New("webhook 不存在")
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, errors.New("生成密钥失败")
	}
	newSecret := hex.EncodeToString(secretBytes)

	if err := l.svcCtx.DB.Model(&record).Update("secret", newSecret).Error; err != nil {
		return nil, errors.New("重置失败")
	}
	return &open_rpc.ResetWebhookSecretRes{Secret: newSecret}, nil
}
