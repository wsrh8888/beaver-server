package message

import (
	"context"
	"errors"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送消息
func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageReq) (resp *types.SendMessageRes, err error) {
	// 1. 从 context 获取 app_id（由中间件注入）
	appID := l.ctx.Value("app_id")
	if appID == nil {
		return nil, errors.New("未认证")
	}

	// 2. 查询 Bot 用户 ID
	var app open_models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ?", appID).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	botUserID := app.BotUserID
	if botUserID == "" {
		return nil, errors.New("Bot 未配置")
	}

	// 3. 调用 chat_rpc 发送消息
	messageId := uuid.New().String()
	_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:    botUserID,
		MessageId: messageId,
		// TODO: 根据 req 构造正确的消息内容
	})
	if err != nil {
		return nil, errors.New("发送消息失败")
	}

	return &types.SendMessageRes{
		MsgID: messageId,
	}, nil
}
