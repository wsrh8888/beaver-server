package logic

import (
	"context"
	"encoding/json"
	"time"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListChatMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListChatMessagesLogic {
	return &ListChatMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListChatMessagesLogic) ListChatMessages(in *chat_rpc.ListChatMessagesReq) (*chat_rpc.ListChatMessagesRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&chat_models.ChatMessage{})
	if in.MessageId != "" {
		db = db.Where("message_id = ?", in.MessageId)
	}
	if in.ConversationId != "" {
		db = db.Where("conversation_id = ?", in.ConversationId)
	}
	if in.SendUserId != "" {
		db = db.Where("send_user_id = ?", in.SendUserId)
	}
	if in.MsgType != 0 {
		db = db.Where("msg_type = ?", in.MsgType)
	}
	if in.Status != 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02 15:04:05", in.StartTime); err == nil {
			db = db.Where("created_at >= ?", startTime)
		}
	}
	if in.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02 15:04:05", in.EndTime); err == nil {
			db = db.Where("created_at <= ?", endTime)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计聊天消息失败: %v", err)
		return nil, err
	}

	var list []chat_models.ChatMessage
	orderClause := "created_at DESC"
	if in.Order == 1 {
		orderClause = "created_at ASC"
	}
	if err := db.Order(orderClause).Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询聊天消息列表失败: %v", err)
		return nil, err
	}

	items := make([]*chat_rpc.ChatMessageItem, 0, len(list))
	for _, m := range list {
		sendUserID := ""
		if m.SendUserID != nil {
			sendUserID = *m.SendUserID
		}
		item := &chat_rpc.ChatMessageItem{
			MessageId:        m.MessageID,
			ConversationId:   m.ConversationID,
			ConversationType: int32(m.ConversationType),
			SendUserId:       sendUserID,
			MsgType:          int32(m.MsgType),
			MsgPreview:       m.MsgPreview,
			Status:           int32(m.Status),
			CreatedAt:        time.Time(m.CreatedAt).Format(time.RFC3339),
			UpdatedAt:        time.Time(m.UpdatedAt).Format(time.RFC3339),
		}
		if in.WithContent && m.Msg != nil {
			if b, err := json.Marshal(m.Msg); err == nil {
				item.MsgContent = string(b)
			}
		}
		items = append(items, item)
	}
	return &chat_rpc.ListChatMessagesRes{Total: total, List: items}, nil
}
