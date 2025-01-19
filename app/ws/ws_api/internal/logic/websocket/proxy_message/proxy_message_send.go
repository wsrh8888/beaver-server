package proxy_message

import (
	"context"
	"encoding/json"
	"net/http"

	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
)

func HandleProxyMessageSend(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, bodyRaw json.RawMessage, ConversationID string) {
	if websocket_utils.IsGroupChat(ConversationID) {

	} else {
		recipientID := websocket_utils.GetRecipientIdFromConversationID(ConversationID, req.UserID)
		websocket_utils.SendMsgByUser(recipientID, req.UserID, "COMMON_PROXY_MESSAGE", type_struct.WsContent{
			Timestamp: 0,
			Data: type_struct.WsData{
				Type:           "transform_websocket_message",
				ConversationID: ConversationID,
				Body:           bodyRaw,
			},
		})
	}
}
