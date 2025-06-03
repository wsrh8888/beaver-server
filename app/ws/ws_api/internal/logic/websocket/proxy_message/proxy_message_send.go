package proxy_message

import (
	"context"
	"encoding/json"
	"net/http"

	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
)

func HandleProxyMessageSend(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, bodyRaw json.RawMessage, ConversationID string) {
	if websocket_utils.IsGroupChat(ConversationID) {

	} else {
		recipientID := websocket_utils.GetRecipientIdFromConversationID(ConversationID, req.UserID)

		content := type_struct.WsContent{
			Timestamp: 0,
			Data: type_struct.WsData{
				Type:           "transform_websocket_message",
				ConversationID: ConversationID,
				Body:           bodyRaw,
			},
		}

		// 分别给接收者和发送者发送消息
		websocket_utils.SendMsgToUser(recipientID, wsCommandConst.FRIEND_OPERATION, content)
		websocket_utils.SendMsgToUser(req.UserID, wsCommandConst.FRIEND_OPERATION, content)
	}
}
