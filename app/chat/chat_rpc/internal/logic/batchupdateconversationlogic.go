package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateConversationLogic {
	return &BatchUpdateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchUpdateConversationLogic) BatchUpdateConversation(in *chat_rpc.BatchUpdateConversationReq) (*chat_rpc.BatchUpdateConversationRes, error) {
	// 开启事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 确保在函数返回时处理事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 批量更新或创建会话记录
	for _, userID := range in.UserIds {
		var userConvo chat_models.ChatUserConversationModel
		err := tx.Where("conversation_id = ? AND user_id = ?", in.ConversationId, userID).First(&userConvo).Error
		if err != nil {
			// 如果记录不存在，创建新记录
			if err := tx.Create(&chat_models.ChatUserConversationModel{
				UserID:         userID,
				ConversationID: in.ConversationId,
				LastMessage:    in.LastMessage,
				IsDeleted:      false,
			}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 如果记录存在，更新记录
			if err := tx.Model(&userConvo).Updates(map[string]interface{}{
				"last_message": in.LastMessage,
				"is_deleted":   false,
			}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &chat_rpc.BatchUpdateConversationRes{
		Success: true,
	}, nil
}
