package pyagent

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

// Client beaver-agent（Python）流式调用。
type Client interface {
	StreamAgentMessage(ctx context.Context, in *StreamAgentMessageReq, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamAgentMessageEvent], error)
}

type client struct {
	cli zrpc.Client
}

func NewClient(cli zrpc.Client) Client {
	return &client{cli: cli}
}

func (c *client) StreamAgentMessage(ctx context.Context, in *StreamAgentMessageReq, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamAgentMessageEvent], error) {
	return NewAgentClient(c.cli.Conn()).StreamAgentMessage(ctx, in, opts...)
}
