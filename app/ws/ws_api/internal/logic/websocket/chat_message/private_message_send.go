package chat_message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	websocket_types "beaver/app/ws/ws_api/internal/logic/websocket/types"
	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/ajax"
	"beaver/common/etcd"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// HandlePrivateMessageSend 处理私聊消息发送
func HandlePrivateMessageSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	r *http.Request,
	conn *websocket.Conn,
	messageId string,
	bodyRaw json.RawMessage,
) {
	fmt.Println("私聊消息开始代理")
	var body type_struct.BodySendMsg
	err := json.Unmarshal(bodyRaw, &body)

	if err != nil {
		fmt.Println("私聊消息解析错误", err.Error())
		return
	}

	apiRequest := websocket_types.SendMsgReq{
		ConversationID: body.ConversationID,
		MessageID:      messageId,
		Msg:            body.Msg,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		fmt.Println("私聊请求数据序列化错误", err)
		return
	}

	addr := etcd.GetServiceAddr(svcCtx.Config.Etcd, "chat_api")
	if addr == "" {
		logx.Error("未匹配到服务")
		return
	}
	apiEndpoint := fmt.Sprintf("http://%s/api/chat/sendMsg", addr)

	sendAjax, err := ajax.ForwardMessage(ajax.ForwardRequest{
		ApiEndpoint: apiEndpoint,
		Method:      "POST",
		UserID:      req.UserID,
		Token:       req.Token,
		Body:        bytes.NewBuffer(requestBody),
	})
	if err != nil {
		fmt.Println("私聊消息发送失败", err)
		return
	}

	// 将 sendAjax 转换为 JSON 格式
	sendAjaxJSON, err := json.Marshal(sendAjax)
	if err != nil {
		fmt.Println("私聊 sendAjax 序列化错误", err)
		return
	}

	// 处理私聊消息转发
	println("当前会话是私聊")
	recipientID := websocket_utils.GetRecipientIdFromConversationID(body.ConversationID, req.UserID)

	// 发送给接收者的消息
	receiverContent := type_struct.WsContent{
		Timestamp: 0,
		MessageID: messageId,
		Data: type_struct.WsData{
			Type:           wsTypeConst.PrivateMessageReceive,
			ConversationID: body.ConversationID,
			Body:           json.RawMessage(sendAjaxJSON),
		},
	}

	// 发送给发送方设备的确认消息（包括发送方自己）
	senderConfirmContent := type_struct.WsContent{
		Timestamp: 0,
		MessageID: messageId,
		Data: type_struct.WsData{
			Type:           wsTypeConst.PrivateMessageSync, // 发送方也用sync类型，表示消息已处理
			ConversationID: body.ConversationID,
			Body:           json.RawMessage(sendAjaxJSON),
		},
	}

	// 1. 发送给接收者的所有设备
	websocket_utils.SendMsgToUser(recipientID, wsCommandConst.CHAT_MESSAGE, receiverContent)

	// 2. 发送给发送者的所有设备（包括发送方设备）
	websocket_utils.SendMsgToUser(req.UserID, wsCommandConst.CHAT_MESSAGE, senderConfirmContent)
}
