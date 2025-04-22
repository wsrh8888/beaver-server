package chat_message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/group/group_rpc/types/group_rpc"
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

func HandleChatMessageSend(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, bodyRaw json.RawMessage) {
	fmt.Println("消息开始代理")
	var body type_struct.BodySendMsg
	err := json.Unmarshal(bodyRaw, &body)

	if err != nil {
		fmt.Println("消息解析错误", err.Error())
		return
	}

	apiRequest := websocket_types.SendMsgReq{
		ConversationID: body.ConversationID,
		Msg:            body.Msg,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		fmt.Println("请求数据序列化错误", err)
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
		fmt.Println("消息发送失败", err)
		return
	}
	// 将 sendAjax 转换为 JSON 格式
	sendAjaxJSON, err := json.Marshal(sendAjax)
	if err != nil {
		fmt.Println("sendAjax 序列化错误", err)
		return
	}
	// 判断是否为群聊，通过 "_" 判断
	if websocket_utils.IsGroupChat(body.ConversationID) {
		println("当前会话是群聊")
		res, _ := svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: body.ConversationID,
		})
		fmt.Println("群聊成员列表", res.Members)
		for _, memberInfo := range res.Members {
			if memberInfo.UserID != req.UserID {
				websocket_utils.SendMsgByUser(memberInfo.UserID, req.UserID, wsCommandConst.CHAT_MESSAGE, type_struct.WsContent{
					Timestamp: 0,
					Data: type_struct.WsData{
						Type: wsTypeConst.GroupMessageSend,
						Body: json.RawMessage(sendAjaxJSON),
					},
				})
			}
		}
	} else {
		println("当前会话是私聊")

		recipientID := websocket_utils.GetRecipientIdFromConversationID(body.ConversationID, req.UserID)
		websocket_utils.SendMsgByUser(recipientID, req.UserID, wsCommandConst.CHAT_MESSAGE, type_struct.WsContent{
			Timestamp: 0,
			Data: type_struct.WsData{
				Type: "private_message_send",
				Body: json.RawMessage(sendAjaxJSON),
			},
		})
	}
}
