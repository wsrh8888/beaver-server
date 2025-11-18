package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearConversationLogic {
	return &ClearConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearConversationLogic) ClearConversation(req *types.ClearConversationReq) (resp *types.ClearConversationRes, err error) {
	if req.ConversationID == "" {
		return nil, errors.New("会话ID不能为空")
	}

	// 逻辑删除会话中的所有消息
	err = l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("conversation_id = ?", req.ConversationID).
		Update("is_deleted", true).Error
	if err != nil {
		logx.Errorf("清空会话消息失败: %v", err)
		return nil, errors.New("清空会话消息失败")
	}

	return &types.ClearConversationRes{}, nil
}
