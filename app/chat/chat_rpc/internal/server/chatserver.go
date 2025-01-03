// Code generated by goctl. DO NOT EDIT.
// Source: chat_rpc.proto

package server

import (
	"context"

	"beaver/app/chat/chat_rpc/internal/logic"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
)

type ChatServer struct {
	svcCtx *svc.ServiceContext
	chat_rpc.UnimplementedChatServer
}

func NewChatServer(svcCtx *svc.ServiceContext) *ChatServer {
	return &ChatServer{
		svcCtx: svcCtx,
	}
}

func (s *ChatServer) SendMsg(ctx context.Context, in *chat_rpc.SendMsgReq) (*chat_rpc.SendMsgRes, error) {
	l := logic.NewSendMsgLogic(ctx, s.svcCtx)
	return l.SendMsg(in)
}