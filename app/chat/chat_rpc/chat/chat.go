// Code generated by goctl. DO NOT EDIT.
// Source: chat_rpc.proto

package chat

import (
	"context"

	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	FileMsg    = chat_rpc.FileMsg
	ImageMsg   = chat_rpc.ImageMsg
	Msg        = chat_rpc.Msg
	SendMsgReq = chat_rpc.SendMsgReq
	SendMsgRes = chat_rpc.SendMsgRes
	Sender     = chat_rpc.Sender
	TextMsg    = chat_rpc.TextMsg
	VideoMsg   = chat_rpc.VideoMsg
	VoiceMsg   = chat_rpc.VoiceMsg

	Chat interface {
		SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRes, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{
		cli: cli,
	}
}

func (m *defaultChat) SendMsg(ctx context.Context, in *SendMsgReq, opts ...grpc.CallOption) (*SendMsgRes, error) {
	client := chat_rpc.NewChatClient(m.cli.Conn())
	return client.SendMsg(ctx, in, opts...)
}
