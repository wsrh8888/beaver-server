package chat_message

import (
	"context"
	"fmt"
	"net/http"

	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/gorilla/websocket"
)

func HandleChatMessageTypes(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	switch content.Data.Type {
	case wsTypeConst.GroupMessageSend:
		HandleGroupMessageSend(ctx, svcCtx, req, r, conn, content.MessageID, content.Data.Body)
	case wsTypeConst.PrivateMessageSend:
		HandlePrivateMessageSend(ctx, svcCtx, req, r, conn, content.MessageID, content.Data.Body)
	default:
		fmt.Println("未支持的消息类型", content.Data.Type)
	}
}
