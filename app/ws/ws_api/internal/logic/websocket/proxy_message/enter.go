package proxy_message

import (
	"context"
	"fmt"
	"net/http"

	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
)

func HandleProxyMessageTypes(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	switch content.Data.Type {
	case "transform_websocket_message":
		HandleProxyMessageSend(ctx, svcCtx, req, r, conn, content.Data.Body, content.Data.ConversationID)
	default:
		fmt.Println("未支持的消息类型3", content.Data.Type)
	}
}
