package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 生成群机器人 Webhook URL（对标钉钉/企业微信）
func NewGenerateWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateWebhookLogic {
	return &GenerateWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateWebhookLogic) GenerateWebhook(req *types.GenerateWebhookReq) (resp *types.GenerateWebhookRes, err error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 2. 生成 Webhook Token
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	// 3. 生成签名密钥
	secret := req.SecretKey
	if secret == "" {
		secretBytes := make([]byte, 16)
		_, _ = rand.Read(secretBytes)
		secret = hex.EncodeToString(secretBytes)
	}

	// 4. 保存 Webhook 配置（简化版，实际应该保存到数据库）
	// TODO: 这里需要根据实际需求调整，目前先返回成功

	// 5. 构造 Webhook URL
	webhookURL := fmt.Sprintf("https://api.beaver.im/webhook/%s/%s", req.AppID, token)

	// 6. 设置过期时间（7天后）
	expireAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	return &types.GenerateWebhookRes{
		WebhookURL: webhookURL,
		Secret:     secret,
		ExpireAt:   expireAt,
	}, nil
}
