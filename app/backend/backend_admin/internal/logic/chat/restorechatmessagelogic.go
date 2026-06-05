package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RestoreChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreChatMessageLogic {
	return &RestoreChatMessageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RestoreChatMessageLogic) RestoreChatMessage(req *types.RestoreChatMessageReq) (resp *types.RestoreChatMessageRes, err error) {
	if req.MessageID == "" {
		return nil, errors.New("消息ID不能为空")
	}

	_, err = l.svcCtx.ChatRpc.UpdateChatMessages(l.ctx, &chat_rpc.UpdateChatMessagesReq{
		MessageIds: []string{req.MessageID},
		Status:     1,
	})
	if err != nil {
		l.Errorf("恢复聊天消息失败: %v", err)
		return nil, err
	}
	return &types.RestoreChatMessageRes{}, nil
}
