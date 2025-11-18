package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationsLogic {
	return &GetUserConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConversationsLogic) GetUserConversations(in *chat_rpc.GetUserConversationsReq) (*chat_rpc.GetUserConversationsRes, error) {
	var userConversations []chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话失败: %v", err)
		return nil, err
	}

	var conversations []*chat_rpc.ConversationItem
	for _, uc := range userConversations {
		// 查询会话类型
		var conversation chat_models.ChatConversationMeta
		err := l.svcCtx.DB.Where("conversation_id = ?", uc.ConversationID).First(&conversation).Error
		if err != nil {
			l.Errorf("查询会话信息失败: %v", err)
			continue
		}

		conversations = append(conversations, &chat_rpc.ConversationItem{
			ConversationId: uc.ConversationID,
			Type:           int32(conversation.Type),
		})
	}

	return &chat_rpc.GetUserConversationsRes{
		Conversations: conversations,
	}, nil
}
