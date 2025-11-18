package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRestoreChatMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量恢复消息
func NewBatchRestoreChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRestoreChatMessagesLogic {
	return &BatchRestoreChatMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRestoreChatMessagesLogic) BatchRestoreChatMessages(req *types.BatchRestoreChatMessagesReq) (resp *types.BatchRestoreChatMessagesRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("消息ID列表不能为空")
	}

	// 批量恢复消息
	err = l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("message_id IN ?", req.Ids).
		Update("is_deleted", false).Error
	if err != nil {
		logx.Errorf("批量恢复聊天消息失败: %v", err)
		return nil, errors.New("批量恢复聊天消息失败")
	}

	return &types.BatchRestoreChatMessagesRes{}, nil
}
