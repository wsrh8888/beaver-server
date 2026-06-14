package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteChatMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchDeleteChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteChatMessagesLogic {
	return &BatchDeleteChatMessagesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *BatchDeleteChatMessagesLogic) BatchDeleteChatMessages(req *types.BatchDeleteChatMessagesReq) (resp *types.BatchDeleteChatMessagesRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("消息ID列表不能为空")
	}

	_, err = l.svcCtx.ChatRpc.UpdateChatMessages(l.ctx, &chat_rpc.UpdateChatMessagesReq{
		MessageIds: req.Ids,
		Status:     chatMessageStatusDeleted,
	})
	if err != nil {
		l.Errorf("批量删除聊天消息失败: %v", err)
		return nil, err
	}
	return &types.BatchDeleteChatMessagesRes{}, nil
}
