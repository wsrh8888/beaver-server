package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取会话数据
func NewGetConversationsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsListByIdsLogic {
	return &GetConversationsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsListByIdsLogic) GetConversationsListByIds(req *types.GetConversationsListByIdsReq) (resp *types.GetConversationsListByIdsRes, err error) {
	// 直接从数据库查询会话完整信息
	var conversations []chat_models.ChatConversationMeta
	err = l.svcCtx.DB.Where("conversation_id IN (?)", req.ConversationIds).Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话信息失败: %v", err)
		return nil, err
	}

	// 转换数据库模型为API响应
	conversationList := make([]types.ConversationById, 0, len(conversations))
	for _, conv := range conversations {
		conversationList = append(conversationList, types.ConversationById{
			ConversationID: conv.ConversationID,
			Type:           conv.Type,
			MaxSeq:         conv.MaxSeq,
			LastMessage:    conv.LastMessage,
			Version:        conv.Version,
			CreateAt:       time.Time(conv.CreatedAt).Unix(),
			UpdateAt:       time.Time(conv.UpdatedAt).Unix(),
		})
	}

	return &types.GetConversationsListByIdsRes{
		Conversations: conversationList,
	}, nil
}
