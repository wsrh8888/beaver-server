package logic

import (
	"context"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	uuidUtil "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBotLogic {
	return &CreateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateBot 创建推送机器人（open 服务负责生成 Token、安全设置等）
func (l *CreateBotLogic) CreateBot(in *open_rpc.CreateBotReq) (*open_rpc.CreateBotRes, error) {
	// 1. 使用传入的机器人用户 ID（由 user_rpc 生成）
	if in.BotId == "" {
		return nil, fmt.Errorf("bot_id 不能为空")
	}

	// 2. 生成 Webhook Token（用于身份验证）
	webhookToken := uuidUtil.NewV4().String()

	// 3. 生成签名密钥（默认启用签名校验）
	signatureSecret := uuidUtil.NewV4().String()

	// 4. 生成 Webhook URL
	webhookURL := fmt.Sprintf("/api/webhook/%s", webhookToken)

	// 5. 创建 open_bots 记录（默认启用签名校验）
	bot := &open_models.OpenBotModel{
		BotID:   in.BotId, // 使用传入的用户 ID
		GroupID: in.GroupId,
		Token:   webhookToken,
		Status:  1,
		Security: open_models.OpenBotSecurity{
			SignatureEnabled:   true, // 默认启用签名校验
			SignatureSecret:    signatureSecret,
			KeywordsEnabled:    false, // 关键词校验默认关闭
			IPWhitelistEnabled: false, // IP白名单默认关闭
		},
	}

	if err := l.svcCtx.DB.Create(bot).Error; err != nil {
		logx.Errorf("创建 Bot 记录失败: %v", err)
		return nil, fmt.Errorf("创建 Bot 记录失败")
	}

	return &open_rpc.CreateBotRes{
		Id:         uint32(bot.ID),
		BotId:      in.BotId,
		WebhookUrl: webhookURL,
		Token:      webhookToken,
	}, nil
}
