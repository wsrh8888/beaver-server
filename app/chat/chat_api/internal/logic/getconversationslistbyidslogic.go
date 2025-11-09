package logic

import (
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"context"
	"errors"
	"time"

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
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	if len(req.ConversationIds) == 0 {
		return &types.GetConversationsListByIdsRes{
			Conversations: []types.ConversationSyncItem{},
		}, nil
	}

	// 查询指定会话ID的会话数据
	var conversations []chat_models.ChatConversationMeta
	err = l.svcCtx.DB.Where("conversation_id IN (?)", req.ConversationIds).Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话数据失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var conversationItems []types.ConversationSyncItem
	for _, conv := range conversations {
		conversationItems = append(conversationItems, types.ConversationSyncItem{
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
		Conversations: conversationItems,
	}, nil
}
