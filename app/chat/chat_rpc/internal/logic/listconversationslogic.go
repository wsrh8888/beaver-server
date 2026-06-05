package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationsLogic {
	return &ListConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListConversationsLogic) ListConversations(in *chat_rpc.ListConversationsReq) (*chat_rpc.ListConversationsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).Where("user_id = ?", in.UserId)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计用户会话失败: %v", err)
		return nil, err
	}

	var userConvos []chat_models.ChatUserConversation
	if err := db.Order("updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&userConvos).Error; err != nil {
		l.Errorf("查询用户会话失败: %v", err)
		return nil, err
	}

	items := make([]*chat_rpc.ConversationDetailItem, 0, len(userConvos))
	for _, uc := range userConvos {
		var meta chat_models.ChatConversationMeta
		if err := l.svcCtx.DB.Where("conversation_id = ?", uc.ConversationID).First(&meta).Error; err != nil {
			continue
		}
		if in.ConversationType > 0 && int32(meta.Type) != in.ConversationType {
			continue
		}

		var msgCount int64
		_ = l.svcCtx.DB.Model(&chat_models.ChatMessage{}).Where("conversation_id = ?", uc.ConversationID).Count(&msgCount).Error

		var participants []chat_models.ChatUserConversation
		_ = l.svcCtx.DB.Where("conversation_id = ?", uc.ConversationID).Find(&participants).Error
		participantIDs := make([]string, 0, len(participants))
		for _, p := range participants {
			participantIDs = append(participantIDs, p.UserID)
		}

		lastTime := ""
		var lastMsg chat_models.ChatMessage
		if err := l.svcCtx.DB.Where("conversation_id = ?", uc.ConversationID).
			Order("created_at DESC").First(&lastMsg).Error; err == nil {
			lastTime = time.Time(lastMsg.CreatedAt).Format(time.RFC3339)
		}

		items = append(items, &chat_rpc.ConversationDetailItem{
			ConversationId:   uc.ConversationID,
			Type:             int32(meta.Type),
			LastMessage:      meta.LastMessage,
			LastMessageTime:  lastTime,
			MessageCount:     msgCount,
			ParticipantIds:   participantIDs,
		})
	}

	return &chat_rpc.ListConversationsRes{Total: total, List: items}, nil
}
