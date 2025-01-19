package websocket_utils

import (
	"encoding/json"
	"fmt"
	"strings"

	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var UserOnlineWsMap = make(map[string]*UserWsInfo)

type UserWsInfo struct {
	WsClientMap map[string]*websocket.Conn // 用户管理的所有 WebSocket 通道
}

// SendMsgByUser 发消息  给指定的用户， 谁发的
func SendMsgByUser(revUserID string, sendUserID string, command string, content type_struct.WsContent) {
	revUser, ok1 := UserOnlineWsMap[revUserID]
	_, ok2 := UserOnlineWsMap[sendUserID]

	if revUserID != sendUserID && ok1 && ok2 {
		jsonContent, _ := json.Marshal(content)
		logx.Info("发送消息给用户数：", len(revUser.WsClientMap), "发送者：", sendUserID, "接收者：", revUserID, "消息内容：", string(jsonContent))
		sendWsMapMsg(revUser.WsClientMap, command, content)
		return
	}
}

// sendWsMapMsg 给一组的 WebSocket 通道发送消息
func sendWsMapMsg(wsMap map[string]*websocket.Conn, command string, content type_struct.WsContent) {
	for _, conn := range wsMap {
		ws_response.WsResponse(conn, command, content)
	}
}

func IsGroupChat(conversationID string) bool {
	return !strings.Contains(conversationID, "_")
}

func GetRecipientIdFromConversationID(conversationID string, userID string) string {
	ids := strings.Split(conversationID, "_")
	if len(ids) != 2 {
		fmt.Println("无效的会话Id：", conversationID)
		return ""
	}
	if ids[0] == userID {
		return ids[1]
	}
	return ids[0]
}
