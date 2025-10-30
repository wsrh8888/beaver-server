package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteChatMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除聊天消息
func NewBatchDeleteChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteChatMessagesLogic {
	return &BatchDeleteChatMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteChatMessagesLogic) BatchDeleteChatMessages(req *types.BatchDeleteChatMessagesReq) (resp *types.BatchDeleteChatMessagesRes, err error) {
	// 批量逻辑删除
	err = l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("message_id IN ?", req.Ids).
		Update("is_deleted", true).Error
	if err != nil {
		logx.Errorf("批量删除聊天消息失败: %v", err)
		return nil, errors.New("批量删除聊天消息失败")
	}

	return &types.BatchDeleteChatMessagesRes{}, nil
}
