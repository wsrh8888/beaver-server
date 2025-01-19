package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/ws/ws_api/internal/logic/websocket/chat_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/proxy_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/webrtc_message"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
)

func HandleWebSocketMessages(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		var wsMessage type_struct.WsMessage
		fmt.Println("收到ws消息", string(p))
		err = json.Unmarshal(p, &wsMessage)
		if err != nil {
			fmt.Println("消息解析错误", err.Error())
			continue
		}

		if wsMessage.Command == "" {
			fmt.Println("command不能为空")
			continue
		}

		switch wsMessage.Command {
		case "COMMON_CHAT_MESSAGE":
			chat_message.HandleChatMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case "COMMON_PROXY_MESSAGE":
			proxy_message.HandleProxyMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case "COMMON_WEBRTC_MESSAGE":
			webrtc_message.HandleWebRTCAnswer(ctx, svcCtx, req, r, conn, wsMessage.Content)
		default:
			fmt.Println("未支持的消息类型", wsMessage.Command)
		}
	}
}
