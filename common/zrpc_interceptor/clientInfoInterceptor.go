package zrpc_interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var clientIP, userId string
	cl := ctx.Value("clientIP")
	if cl != nil {
		clientIP = cl.(string)
	}
	uid := ctx.Value("userId")
	if uid != nil {
		userId = uid.(string)
	}

	md := metadata.New(map[string]string{"clientIP": clientIP, "userId": userId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
