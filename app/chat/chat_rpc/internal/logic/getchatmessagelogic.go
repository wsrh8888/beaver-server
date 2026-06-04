package logic

import (
	"context"
	"errors"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageLogic {
	return &GetChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatMessageLogic) GetChatMessage(in *chat_rpc.GetChatMessageReq) (*chat_rpc.GetChatMessageRes, error) {
	if in.ConversationId == "" || in.MessageId == "" {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	var row chat_models.ChatMessage
	err := l.svcCtx.DB.Where("conversation_id = ? AND message_id = ?", in.ConversationId, in.MessageId).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}
	if err != nil {
		return nil, err
	}
	if row.Msg == nil {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	sl := NewSendMsgLogic(l.ctx, l.svcCtx)
	protoMsg, err := sl.convertCtypeMsgToGrpcMsg(*row.Msg)
	if err != nil || protoMsg == nil {
		return &chat_rpc.GetChatMessageRes{Found: false}, nil
	}

	return &chat_rpc.GetChatMessageRes{
		Found: true,
		Msg:   protoMsg,
	}, nil
}
