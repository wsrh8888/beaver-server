package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWebhookLogic {
	return &CreateWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWebhookLogic) CreateWebhook(in *open_rpc.CreateWebhookReq) (*open_rpc.CreateWebhookRes, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.New("生成 token 失败")
	}
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, errors.New("生成 secret 失败")
	}

	token := hex.EncodeToString(tokenBytes)
	secret := hex.EncodeToString(secretBytes)
	botUserID := "nbot_" + token[:32]

	record := open_models.OpenGroupBotModel{
		Token:     token,
		Secret:    secret,
		AppID:     "GROUP_NOTIFICATION",
		GroupID:   in.GroupId,
		BotUserID: botUserID,
		Status:    1,
	}
	if err := l.svcCtx.DB.Create(&record).Error; err != nil {
		return nil, errors.New("创建 webhook 失败")
	}

	webhookURL := fmt.Sprintf("%s/api/open/v1/robot/send?access_token=%s", l.svcCtx.Config.ApiBaseUrl, token)
	return &open_rpc.CreateWebhookRes{
		Id:         uint32(record.ID),
		BotUserId:  botUserID,
		WebhookUrl: webhookURL,
		Secret:     secret,
	}, nil
}
