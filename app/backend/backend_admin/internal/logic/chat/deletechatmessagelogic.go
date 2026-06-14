package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteChatMessageLogic {
	return &DeleteChatMessageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteChatMessageLogic) DeleteChatMessage(req *types.DeleteChatMessageReq) (resp *types.DeleteChatMessageRes, err error) {
	if req.MessageID == "" {
		return nil, errors.New("消息ID不能为空")
	}

	_, err = l.svcCtx.ChatRpc.UpdateChatMessages(l.ctx, &chat_rpc.UpdateChatMessagesReq{
		MessageIds: []string{req.MessageID},
		Status:     chatMessageStatusDeleted,
	})
	if err != nil {
		l.Errorf("删除聊天消息失败: %v", err)
		return nil, err
	}
	return &types.DeleteChatMessageRes{}, nil
}
