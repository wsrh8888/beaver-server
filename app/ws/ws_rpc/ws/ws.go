// Code generated by goctl. DO NOT EDIT.
// Source: ws_rpc.proto

package ws

import (
	"context"

	"beaver/app/ws/ws_rpc/types/ws_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	HandleWebSocketRequest   = ws_rpc.HandleWebSocketRequest
	HandleWebSocketResponse  = ws_rpc.HandleWebSocketResponse
	SendProxyMessageRequest  = ws_rpc.SendProxyMessageRequest
	SendProxyMessageResponse = ws_rpc.SendProxyMessageResponse

	Ws interface {
		HandleWebSocket(ctx context.Context, in *HandleWebSocketRequest, opts ...grpc.CallOption) (*HandleWebSocketResponse, error)
		SendProxyMessage(ctx context.Context, in *SendProxyMessageRequest, opts ...grpc.CallOption) (*SendProxyMessageResponse, error)
	}

	defaultWs struct {
		cli zrpc.Client
	}
)

func NewWs(cli zrpc.Client) Ws {
	return &defaultWs{
		cli: cli,
	}
}

func (m *defaultWs) HandleWebSocket(ctx context.Context, in *HandleWebSocketRequest, opts ...grpc.CallOption) (*HandleWebSocketResponse, error) {
	client := ws_rpc.NewWsClient(m.cli.Conn())
	return client.HandleWebSocket(ctx, in, opts...)
}

func (m *defaultWs) SendProxyMessage(ctx context.Context, in *SendProxyMessageRequest, opts ...grpc.CallOption) (*SendProxyMessageResponse, error) {
	client := ws_rpc.NewWsClient(m.cli.Conn())
	return client.SendProxyMessage(ctx, in, opts...)
}
