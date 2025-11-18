package chat_message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"

	"github.com/gorilla/websocket"
)

// HandleGroupMessageSend 处理群聊消息发送
func HandleGroupMessageSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	r *http.Request,
	conn *websocket.Conn,
	messageId string,
	bodyRaw json.RawMessage,
) {
	fmt.Println("群聊消息开始代理")

	var body type_struct.BodySendMsg
	err := json.Unmarshal(bodyRaw, &body)
	if err != nil {
		fmt.Println("群聊消息解析错误", err.Error())
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
	_, err = svcCtx.ChatRpc.SendMsg(ctx, rpcReq)
	if err != nil {
		fmt.Println("群聊消息发送失败", err)
		return
	}

	// // 构建响应数据
	// responseJSON, err := buildResponseData(rpcResp, body.Msg)
	// if err != nil {
	// 	fmt.Println("构建响应数据失败", err)
	// 	return
	// }

	// // 处理群聊消息转发
	// println("当前会话是群聊")
	// res, err := svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
	// 	GroupID: body.ConversationID,
	// })
	// if err != nil {
	// 	fmt.Println("获取群聊成员列表失败", err)
	// 	return
	// }

	// fmt.Println("群聊成员列表", res.Members)

	// // 给所有群成员发送消息
	// for _, memberInfo := range res.Members {
	// 	websocket_utils.SendMsgToUser(memberInfo.UserID, wsCommandConst.CHAT_MESSAGE, type_struct.WsContent{
	// 		Timestamp: 0,
	// 		MessageID: messageId,
	// 		Data: type_struct.WsData{
	// 			Type:           wsTypeConst.ChatConversationMessageReceive,
	// 			ConversationID: body.ConversationID,
	// 			Body:           json.RawMessage(responseJSON),
	// 		},
	// 	})
	// }
}
