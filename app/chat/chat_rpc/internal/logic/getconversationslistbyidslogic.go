package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsListByIdsLogic {
	return &GetConversationsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationsListByIdsLogic) GetConversationsListByIds(in *chat_rpc.GetConversationsListByIdsReq) (*chat_rpc.GetConversationsListByIdsRes, error) {
	// 构建查询条件
	query := l.svcCtx.DB.Where("conversation_id IN (?)", in.ConversationIds)

	// 如果提供了since参数，只返回版本有变更的记录
	if in.Since > 0 {
		query = query.Where("version >= ?", in.Since)
	}

	// 查询指定会话的完整信息
	var conversations []chat_models.ChatConversationMeta
	err := query.Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话信息失败: %v", err)
		return nil, err
	}

	var conversationList []*chat_rpc.ConversationListById
	for _, conv := range conversations {
		conversationList = append(conversationList, &chat_rpc.ConversationListById{
			ConversationId: conv.ConversationID,
			Type:           int32(conv.Type),
			Seq:            conv.MaxSeq,
			Version:        conv.Version,
		})
	}

	return &chat_rpc.GetConversationsListByIdsRes{
		Conversations: conversationList,
	}, nil
}
