package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"beaver/app/user/user_models"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/ajax"
	"beaver/common/etcd"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebsocketLogic {
	return &ChatWebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type UserWsInfo struct {
	UserInfo    user_models.UserModel
	WsClientMap map[string]*websocket.Conn //用户管理的所有ws通道
}

var UserOnlineWsMap = map[string]*UserWsInfo{}

type SendMsgReq struct {
	UserId         string          `header:"Beaver-User-Id"`
	ConversationId string          `json:"conversationId"`
	Msg            json.RawMessage `json:"msg"`
}

func (l *ChatWebsocketLogic) ChatWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		return
	}
	defer cleanupConnection(req.UserId, conn)

	fmt.Println("用户上线", req.UserId)
	userWsInfo, ok := UserOnlineWsMap[req.UserId]
	addr := conn.RemoteAddr().String()

	if !ok {
		userWsInfo = &UserWsInfo{
			WsClientMap: map[string]*websocket.Conn{
				addr: conn,
			},
		}
		UserOnlineWsMap[req.UserId] = userWsInfo
	} else {
		// 则替换当前连接
		userWsInfo.WsClientMap[addr] = conn
	}

	l.handleWebSocketMessages(req, r, conn)
	return &types.WsRes{}, nil
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func cleanupConnection(userId string, conn *websocket.Conn) {
	conn.Close()
	addr := conn.RemoteAddr().String()

	userWsInfo, ok := UserOnlineWsMap[userId]
	if ok {
		delete(userWsInfo.WsClientMap, addr)
	}
	if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
		//如果都退出了，就下线
		delete(UserOnlineWsMap, userId)
	}

}

func (l *ChatWebsocketLogic) handleWebSocketMessages(req *types.WsReq, r *http.Request, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// err := conn.WriteMessage(websocket.PingMessage, []byte{})
			err := conn.WriteMessage(websocket.TextMessage, []byte{})
			if err != nil {
				log.Println("WriteMessage error (ping):", err)
				return
			}
		}
	}()

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
			l.handleCommonChatMessage(req, r, conn, wsMessage.Content)
		case "COMMON_PROXY_MESSAGE":
			handleProxyMessage(req, r, conn, wsMessage.Content)
		default:
			fmt.Println("未支持的消息类型", wsMessage.Command)
		}
	}
}

func handleProxyMessage(req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	switch content.Data.Type {
	case "webrtc_screen_message":
		handleScreenMessageSend(req, r, conn, content.Data.Body)
	default:
		fmt.Println("未支持的消息类型", content.Data.Type)
	}
}

func (l *ChatWebsocketLogic) handleCommonChatMessage(req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	switch content.Data.Type {
	case "chat_message_send":
		l.handleChatMessageSend(req, r, conn, content.Data.Body)
	default:
		fmt.Println("未支持的消息类型", content.Data.Type)
	}
}

func handleScreenMessageSend(req *types.WsReq, r *http.Request, conn *websocket.Conn, bodyRaw json.RawMessage) {

	var body type_struct.BodyProxyMsg
	err := json.Unmarshal(bodyRaw, &body)
	if err != nil {
		fmt.Println("消息解析错误", err.Error())
		return
	}
	if isGroupChat(body.ConversationId) {

	} else {
		recipientId := getRecipientIdFromConversationId(body.ConversationId, req.UserId)
		SendMsgByUser(recipientId, req.UserId, "COMMON_PROXY_MESSAGE", type_struct.WsContent{
			Timestamp: 0,
			Data: type_struct.WsData{
				Type: "webrtc_screen_message",
				Body: json.RawMessage(bodyRaw),
			},
		})
	}
}

func (l *ChatWebsocketLogic) handleChatMessageSend(req *types.WsReq, r *http.Request, conn *websocket.Conn, bodyRaw json.RawMessage) {
	fmt.Println("消息开始代理")
	var body type_struct.BodySendMsg
	err := json.Unmarshal(bodyRaw, &body)
	if err != nil {
		fmt.Println("消息解析错误", err.Error())
		return
	}

	apiRequest := SendMsgReq{
		ConversationId: body.ConversationId,
		Msg:            body.Msg,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		fmt.Println("请求数据序列化错误", err)
		return
	}

	addr := etcd.GetServiceAddr(l.svcCtx.Config.Etcd, "chat_api")
	if addr == "" {
		logx.Error("未匹配到服务")
		return
	}
	apiEndpoint := fmt.Sprintf("http://%s/api/chat/sendMsg", addr)

	sendAjax, err := ajax.ForwardMessage(ajax.ForwardRequest{
		ApiEndpoint: apiEndpoint,
		Method:      "POST",
		UserId:      req.UserId,
		Token:       req.Token,
		Body:        bytes.NewBuffer(requestBody),
	})
	if err != nil {
		fmt.Println("消息发送失败", err)
		return
	}
	// 将 sendAjax 转换为 JSON 格式
	sendAjaxJSON, err := json.Marshal(sendAjax)
	if err != nil {
		fmt.Println("sendAjax 序列化错误", err)
		return
	}

	recipientId := getRecipientIdFromConversationId(body.ConversationId, req.UserId)
	SendMsgByUser(recipientId, req.UserId, "COMMON_CHAT_MESSAGE", type_struct.WsContent{
		Timestamp: 0,
		Data: type_struct.WsData{
			Type: "private_message_send",
			Body: json.RawMessage(sendAjaxJSON),
		},
	})

}

/**
 * @description: 发消息  给谁发， 谁发的
 */
func SendMsgByUser(revUserId string, sendUserId string, command string, content type_struct.WsContent) {
	revUser, ok1 := UserOnlineWsMap[revUserId]
	_, ok2 := UserOnlineWsMap[sendUserId]

	if revUserId != sendUserId && ok1 && ok2 {
		jsonContent, _ := json.Marshal(content)
		logx.Info("发送消息给用户数：", len(revUser.WsClientMap), "发送者：", sendUserId, "接收者：", revUserId, "消息内容：", string(jsonContent))
		sendWsMapMsg(revUser.WsClientMap, command, content)
		return
	}
}

/**
 * @description: 给一组的ws通道发送消息
 */
func sendWsMapMsg(wsMap map[string]*websocket.Conn, command string, content type_struct.WsContent) {
	for _, conn := range wsMap {
		ws_response.WsResponse(conn, command, content)
	}
}

func isGroupChat(conversationId string) bool {
	return !strings.Contains(conversationId, "_")
}

func getRecipientIdFromConversationId(conversationId string, userId string) string {
	ids := strings.Split(conversationId, "_")
	if len(ids) != 2 {
		fmt.Println("无效的会话Id：", conversationId)
		return ""
	}
	if ids[0] == userId {
		return ids[1]
	}
	return ids[0]
}
