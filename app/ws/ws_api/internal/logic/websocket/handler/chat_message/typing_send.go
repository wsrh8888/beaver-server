package chat_message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/chat/chat_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

func HandleTypingSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	_ *http.Request,
	_ *ws_conn.Client,
	bodyRaw json.RawMessage,
) error {
	var body type_struct.BodyTyping
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		return fmt.Errorf("typing 消息格式错误: %w", err)
	}
	if body.ConversationID == "" {
		return fmt.Errorf("conversationId 不能为空")
	}

	peerIDs, err := getTypingPeerIDs(ctx, svcCtx, req.UserID, body.ConversationID)
	if err != nil {
		return err
	}

	for _, peerID := range peerIDs {
		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     wsTypeConst.TypingReceive,
			"senderId": req.UserID,
			"targetId": peerID,
			"body": map[string]interface{}{
				"conversationId": body.ConversationID,
				"userId":         req.UserID,
				"isTyping":       body.IsTyping,
			},
			"conversationId": body.ConversationID,
		}
		if err := svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logx.Errorf("推送 typing 通知失败: sender=%s, target=%s, error=%v", req.UserID, peerID, err)
		}
	}
	return nil
}

func getTypingPeerIDs(ctx context.Context, svcCtx *svc.ServiceContext, currentUserID, conversationID string) ([]string, error) {
	convType, userIDs := conversation.ParseConversationWithType(conversationID)
	if convType == 1 {
		var peers []string
		for _, uid := range userIDs {
			if uid != currentUserID {
				peers = append(peers, uid)
			}
		}
		if len(peers) == 0 {
			return nil, fmt.Errorf("无效的私聊会话")
		}
		return peers, nil
	}

	groupID := conversation.GetTargetIDByConversation(conversationID, currentUserID)
	membersRes, err := svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, fmt.Errorf("查询群成员失败: %w", err)
	}

	peers := make([]string, 0, len(membersRes.Members))
	for _, m := range membersRes.Members {
		if m.UserID != currentUserID {
			peers = append(peers, m.UserID)
		}
	}
	if len(peers) == 0 {
		var userConversations []chat_models.ChatUserConversation
		if err := svcCtx.DB.Where("conversation_id = ? AND user_id <> ?", conversationID, currentUserID).
			Find(&userConversations).Error; err != nil {
			return nil, err
		}
		for _, uc := range userConversations {
			peers = append(peers, uc.UserID)
		}
	}
	return peers, nil
}
