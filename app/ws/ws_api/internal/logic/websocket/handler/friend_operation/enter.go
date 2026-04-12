package friend_operation

import (
	"context"
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
)

// Handle 处理好友操作类消息（客户端→服务端）
// 例如：发起好友申请、接受/拒绝好友申请等，需写入 DB
func Handle(_ context.Context, _ *svc.ServiceContext, _ *types.WsReq, _ *http.Request, _ *ws_conn.Client, _ type_struct.WsContent) error {
	// TODO: 实现好友操作逻辑（写 DB）
	return nil
}
