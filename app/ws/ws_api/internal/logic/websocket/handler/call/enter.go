package call

import (
	"context"
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
)

// Handle 处理音视频通话信令（客户端→服务端，转发给对端）
func Handle(_ context.Context, _ *svc.ServiceContext, _ *types.WsReq, _ *http.Request, _ *ws_conn.Client, _ type_struct.WsContent) error {
	// TODO: 实现通话信令转发逻辑
	return nil
}
