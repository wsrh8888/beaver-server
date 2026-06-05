package logic

import (
	"context"
	"time"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChatMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChatMessagesLogic {
	return &UpdateChatMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateChatMessagesLogic) UpdateChatMessages(in *chat_rpc.UpdateChatMessagesReq) (*chat_rpc.UpdateChatMessagesRes, error) {
	db := l.svcCtx.DB.Model(&chat_models.ChatMessage{})
	if len(in.MessageIds) > 0 {
		db = db.Where("message_id IN ?", in.MessageIds)
	}
	if in.ConversationId != "" {
		db = db.Where("conversation_id = ?", in.ConversationId)
	}
	if in.MsgType != 0 {
		db = db.Where("msg_type = ?", in.MsgType)
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

	var count int64
	if err := db.Count(&count).Error; err != nil {
		l.Errorf("统计消息失败: %v", err)
		return nil, err
	}
	if err := db.Update("status", in.Status).Error; err != nil {
		l.Errorf("更新消息状态失败: %v", err)
		return nil, err
	}
	return &chat_rpc.UpdateChatMessagesRes{AffectedCount: count}, nil
}
