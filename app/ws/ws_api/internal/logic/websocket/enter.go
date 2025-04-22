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
	"beaver/common/wsEnum/wsCommandConst"

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

		switch wsCommandConst.Command(wsMessage.Command) {
		case wsCommandConst.CHAT_MESSAGE:
			chat_message.HandleChatMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case wsCommandConst.FRIEND_OPERATION:
			proxy_message.HandleProxyMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case wsCommandConst.GROUP_OPERATION:
			webrtc_message.HandleWebRTCAnswer(ctx, svcCtx, req, r, conn, wsMessage.Content)
		default:
			fmt.Println("未支持的消息类型1", wsMessage.Command)
		}
	}
}
