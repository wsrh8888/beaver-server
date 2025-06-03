package grpcMiddleware

import (
	"beaver/common/middle/utils"
	"context"
	"time"

	"google.golang.org/grpc"
)

// RequestLogInterceptor gRPC 请求日志拦截器
func RequestLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()

	// 调用下一个处理器
	resp, err := handler(ctx, req)
	// 记录请求信息
	utils.LogRequest("gRPC", info.FullMethod, req, resp, err, startTime)

	return resp, err
}
