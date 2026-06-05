package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearConversationLogic {
	return &ClearConversationLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ClearConversationLogic) ClearConversation(req *types.ClearConversationReq) (resp *types.ClearConversationRes, err error) {
	if req.ConversationID == "" {
		return nil, errors.New("会话ID不能为空")
	}

	_, err = l.svcCtx.ChatRpc.UpdateChatMessages(l.ctx, &chat_rpc.UpdateChatMessagesReq{
		ConversationId: req.ConversationID,
		Status:         chatMessageStatusDeleted,
	})
	if err != nil {
		l.Errorf("清空会话消息失败: %v", err)
		return nil, err
	}
	return &types.ClearConversationRes{}, nil
}
