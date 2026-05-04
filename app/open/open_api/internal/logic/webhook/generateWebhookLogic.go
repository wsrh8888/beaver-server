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
	models "beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateWebhookLogic {
	return &GenerateWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateWebhookLogic) GenerateWebhook(req *types.GenerateWebhookReq) (resp *types.GenerateWebhookRes, err error) {
	// 1. 验证应用是否存在
	var app models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	// 2. 生成随机 token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	// 3. 生成签名密钥
	secretBytes := make([]byte, 16)
	rand.Read(secretBytes)
	secret := hex.EncodeToString(secretBytes)

	// 4. 构建 Webhook URL (暂时使用占位符,实际应该从配置或环境变量获取)
	webhookURL := fmt.Sprintf("https://your-domain.com/api/open/v1/webhook/send?token=%s", token)

	// 5. 保存到数据库
	webhookConfig := models.OpenWebhookConfig{
		AppID:     req.AppID,
		EventType: "group_message", // 群消息事件
		TargetURL: webhookURL,
		Secret:    secret,
		Status:    1, // 启用
	}

	err = l.svcCtx.DB.Create(&webhookConfig).Error
	if err != nil {
		logx.Errorf("创建 Webhook 配置失败: %v", err)
		return nil, errors.New("生成 Webhook 失败")
	}

	logx.Infof("生成 Webhook 成功: app_id=%s", req.AppID)

	return &types.GenerateWebhookRes{
		WebhookURL: webhookURL,
		Secret:     secret,
		ExpireAt:   time.Now().Add(365 * 24 * time.Hour).Unix(), // 1年后过期
	}, nil
}
