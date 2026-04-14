package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/logic/websocket/handler/chat_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func HandleWebSocketMessages(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, client *ws_conn.Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket连接异常关闭, 用户: %s, 错误: %v\n", req.UserID, err)
			} else {
				fmt.Printf("WebSocket连接正常关闭, 用户: %s\n", req.UserID)
			}
			break
		}

		var wsMessage type_struct.WsMessage
		if err = json.Unmarshal(p, &wsMessage); err != nil {
			fmt.Printf("消息解析错误, 用户: %s, 错误: %s\n", req.UserID, err.Error())
			continue
		}

		if wsMessage.Command == "" {
			continue
		}

		cmd := wsCommandConst.Command(wsMessage.Command)

		// 控制帧：PING/PONG 直接处理，不发 ACK
		switch cmd {
		case wsCommandConst.PING:
			fmt.Printf("收到 PING: 用户: %s, 原始时间戳: %d\n", req.UserID, wsMessage.Content.Timestamp)
			heartbeat.HandleClientPing(client, wsMessage.Content.Timestamp)
			continue
		case wsCommandConst.PONG:
			// 收到客户端对服务端 PING 的回复，无需处理
			fmt.Printf("收到 PONG: 用户: %s, 时间戳: %d\n", req.UserID, wsMessage.Content.Timestamp)
			continue
		case wsCommandConst.USER_PROFILE, wsCommandConst.NOTIFICATION, wsCommandConst.EMOJI:
			// 仅服务端推送，客户端不应发送
			fmt.Printf("客户端不应发送此命令, 用户: %s, 命令: %s\n", req.UserID, cmd)
			continue
		}

		// 业务命令：立即发送 ACK（表示服务端已收到），再处理
		msgId := wsMessage.Content.MessageID
		client.SafeSendControl(type_struct.WsControlFrame{
			Command:   wsCommandConst.ACK,
			MessageID: msgId,
		})

		var handlerErr error
		switch cmd {
		case wsCommandConst.CHAT_MESSAGE:
			handlerErr = chat_message.Handle(ctx, svcCtx, req, r, client, wsMessage.Content)
		default:
			fmt.Printf("未支持的命令类型, 用户: %s, 命令: %s\n", req.UserID, wsMessage.Command)
		}

		if handlerErr != nil {
			logx.Errorf("处理命令失败, 用户: %s, 命令: %s, 错误: %v", req.UserID, cmd, handlerErr)
		}
	}
}
