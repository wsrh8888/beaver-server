package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/ws/ws_api/internal/logic/websocket/chat_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	"beaver/app/ws/ws_api/internal/logic/websocket/proxy_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/webrtc_message"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
)

func HandleWebSocketMessages(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, heartbeatManager *heartbeat.Manager) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			// 检查是否是正常关闭或网络错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket连接异常关闭, 用户: %s, 错误: %v\n", req.UserID, err)
			} else {
				fmt.Printf("WebSocket连接正常关闭, 用户: %s\n", req.UserID)
			}
			break
		}

		var wsMessage type_struct.WsMessage
		fmt.Printf("收到ws消息, 用户: %s, 内容: %s\n", req.UserID, string(p))

		err = json.Unmarshal(p, &wsMessage)
		if err != nil {
			fmt.Printf("消息解析错误, 用户: %s, 错误: %s, 原始消息: %s\n", req.UserID, err.Error(), string(p))
			continue
		}

		if wsMessage.Command == "" {
			fmt.Printf("command不能为空, 用户: %s, 消息: %s\n", req.UserID, string(p))
			continue
		}

		switch wsCommandConst.Command(wsMessage.Command) {
		case wsCommandConst.CHAT_MESSAGE:
			chat_message.HandleChatMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case wsCommandConst.FRIEND_OPERATION:
			proxy_message.HandleProxyMessageTypes(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case wsCommandConst.GROUP_OPERATION:
			webrtc_message.HandleWebRTCAnswer(ctx, svcCtx, req, r, conn, wsMessage.Content)
		case wsCommandConst.HEARTBEAT:
			// 添加调试信息
			fmt.Printf("处理心跳消息, 用户: %s, 命令: %s\n", req.UserID, wsMessage.Command)
			// 使用心跳管理器处理心跳
			heartbeatManager.HandleClientHeartbeat(wsMessage.Content)
		default:
			fmt.Printf("未支持的消息类型, 用户: %s, 命令: %s, 命令类型: %T\n", req.UserID, wsMessage.Command, wsMessage.Command)
			fmt.Printf("HEARTBEAT常量值: %s, 命令值: %s, 是否相等: %v\n", wsCommandConst.HEARTBEAT, wsMessage.Command, wsCommandConst.HEARTBEAT == wsCommandConst.Command(wsMessage.Command))
		}
	}
}
