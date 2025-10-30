package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationLogic {
	return &UpdateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConversationLogic) UpdateConversation(in *chat_rpc.UpdateConversationReq) (*chat_rpc.UpdateConversationRes, error) {
	var userConvo chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", in.ConversationId, in.UserId).First(&userConvo).Error
	if err != nil {
		// 如果记录不存在，创建新记录
		if err := l.svcCtx.DB.Create(&chat_models.ChatUserConversation{
			UserID:         in.UserId,
			ConversationID: in.ConversationId,
			LastMessage:    in.LastMessage,
			IsPinned:       in.IsPinned,
			IsDeleted:      in.IsDeleted,
		}).Error; err != nil {
			return nil, err
		}
	} else {
		// 如果记录存在，更新记录
		updates := map[string]interface{}{
			"is_deleted": in.IsDeleted,
		}
		if in.LastMessage != "" {
			updates["last_message"] = in.LastMessage
		}
		if in.IsPinned != userConvo.IsPinned {
			updates["is_pinned"] = in.IsPinned
		}

		if err := l.svcCtx.DB.Model(&userConvo).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &chat_rpc.UpdateConversationRes{
		Success: true,
	}, nil
}
