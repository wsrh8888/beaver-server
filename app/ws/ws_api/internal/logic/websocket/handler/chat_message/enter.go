package chat_message

import (
	"context"
	"fmt"
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsTypeConst"
)

func Handle(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, client *ws_conn.Client, content type_struct.WsContent) error {
	switch content.Data.Type {
	case wsTypeConst.GroupMessageSend:
		return HandleGroupMessageSend(ctx, svcCtx, req, r, client, content.MessageID, content.Data.Body)
	case wsTypeConst.PrivateMessageSend:
		return HandlePrivateMessageSend(ctx, svcCtx, req, r, client, content.MessageID, content.Data.Body)
	default:
		fmt.Println("未支持的消息类型:", content.Data.Type)
		return fmt.Errorf("unsupported message type: %s", content.Data.Type)
	}
}
