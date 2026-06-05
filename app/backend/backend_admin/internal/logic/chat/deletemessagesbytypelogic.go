package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessagesByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMessagesByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessagesByTypeLogic {
	return &DeleteMessagesByTypeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteMessagesByTypeLogic) DeleteMessagesByType(req *types.DeleteMessagesByTypeReq) (resp *types.DeleteMessagesByTypeRes, err error) {
	if req.MsgType == 0 {
		return nil, errors.New("消息类型不能为空")
	}

	rpcRes, err := l.svcCtx.ChatRpc.UpdateChatMessages(l.ctx, &chat_rpc.UpdateChatMessagesReq{
		ConversationId: req.ConversationID,
		MsgType:        int32(req.MsgType),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Status:         chatMessageStatusDeleted,
	})
	if err != nil {
		l.Errorf("按类型删除消息失败: %v", err)
		return nil, err
	}
	return &types.DeleteMessagesByTypeRes{DeletedCount: rpcRes.AffectedCount}, nil
}
