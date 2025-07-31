package chat_message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	websocket_utils "beaver/app/ws/ws_api/internal/logic/websocket/utils"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/gorilla/websocket"
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

	// 将原始消息转换为RPC消息格式
	rpcMsg, err := convertToRpcMsg(body.Msg)
	if err != nil {
		fmt.Println("消息格式转换错误", err.Error())
		return
	}

	// 构建RPC请求
	rpcReq := &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		ConversationId: body.ConversationID,
		MessageId:      messageId,
		Msg:            rpcMsg,
	}

	// 调用RPC服务
	rpcResp, err := svcCtx.ChatRpc.SendMsg(ctx, rpcReq)
	if err != nil {
		fmt.Println("私聊消息发送失败", err)
		return
	}

	// 构建响应数据
	responseJSON, err := buildResponseData(rpcResp, body.Msg)
	if err != nil {
		fmt.Println("构建响应数据失败", err)
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
			Body:           json.RawMessage(responseJSON),
		},
	}

	// 发送给发送方设备的确认消息（包括发送方自己）
	senderConfirmContent := type_struct.WsContent{
		Timestamp: 0,
		MessageID: messageId,
		Data: type_struct.WsData{
			Type:           wsTypeConst.PrivateMessageSync, // 发送方也用sync类型，表示消息已处理
			ConversationID: body.ConversationID,
			Body:           json.RawMessage(responseJSON),
		},
	}

	// 1. 发送给接收者的所有设备
	websocket_utils.SendMsgToUser(recipientID, wsCommandConst.CHAT_MESSAGE, receiverContent)

	// 2. 发送给发送者的所有设备（包括发送方设备）
	websocket_utils.SendMsgToUser(req.UserID, wsCommandConst.CHAT_MESSAGE, senderConfirmContent)
}
