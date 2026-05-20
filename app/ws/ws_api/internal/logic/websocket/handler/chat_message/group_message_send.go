package chat_message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
)

func HandleGroupMessageSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	r *http.Request,
	client *ws_conn.Client,
	messageId string,
	bodyRaw json.RawMessage,
) error {
	var body type_struct.BodySendMsg
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		fmt.Println("群聊消息解析错误", err.Error())
		return fmt.Errorf("消息格式错误: %w", err)
	}

	rpcMsg, err := convertToRpcMsg(body.Msg)
	if err != nil {
		fmt.Println("消息格式转换错误", err.Error())
		return fmt.Errorf("消息内容错误: %w", err)
	}

	_, err = svcCtx.ChatRpc.SendMsg(ctx, &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		ConversationId: body.ConversationID,
		MessageId:      messageId,
		Msg:            rpcMsg,
	})
	if err != nil {
		fmt.Println("群聊消息发送失败", err)
		return err
	}

	// Bot 流式 Webhook（bot_streaming.go）暂禁用
	// go func() {
	// 	var msgContent string
	// 	var msgMap map[string]interface{}
	// 	if err := json.Unmarshal(body.Msg, &msgMap); err == nil {
	// 		if content, ok := msgMap["content"].(string); ok {
	// 			msgContent = content
	// 		}
	// 	}
	// 	handleBotStreaming(ctx, svcCtx.DB, client, body.ConversationID, msgContent, req.UserID)
	// }()

	return nil
}
