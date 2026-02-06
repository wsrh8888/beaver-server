package webrtc_message

import (
	"beaver/app/call/call_models"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

func HandleWebRTCAnswer(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, conn *websocket.Conn, content type_struct.WsContent) {
	// 解析消息内容，检查是否有Type字段
	var msgData map[string]interface{}
	if err := json.Unmarshal(content.Data.Body, &msgData); err != nil {
		logx.Errorf("解析消息内容失败: %v", err)
		return
	}

	msgType, ok := msgData["type"].(float64)
	if !ok {
		logx.Error("消息缺少type字段")
		return
	}

	switch int(msgType) {
	case 6: // 语音通话
		handleCallMessage(ctx, svcCtx, req, content, "voice")
	case 7: // 视频通话
		handleCallMessage(ctx, svcCtx, req, content, "video")
	default:
		fmt.Println("未支持的通话消息类型", msgType)
	}
}

// handleCallMessage 处理通话相关消息
func handleCallMessage(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, content type_struct.WsContent, callType string) {
	// 解析消息内容
	var msgData map[string]interface{}
	if err := json.Unmarshal(content.Data.Body, &msgData); err != nil {
		logx.Errorf("解析通话消息失败: %v", err)
		return
	}

	action, ok := msgData["action"].(string)
	if !ok {
		logx.Error("通话消息缺少action字段")
		return
	}

	roomID, ok := msgData["roomId"].(string)
	if !ok {
		logx.Error("通话消息缺少roomId字段")
		return
	}

	switch action {
	case "join":
		// 用户加入通话房间
		handleJoinCall(ctx, svcCtx, req.UserID, roomID, callType)

	case "leave":
		// 用户离开通话房间
		handleLeaveCall(ctx, svcCtx, req.UserID, roomID)

	case "heartbeat":
		// 通话心跳，更新用户状态
		handleCallHeartbeat(ctx, svcCtx, req.UserID, roomID)

	default:
		logx.Errorf("未知的通话动作: %s", action)
	}
}

// handleJoinCall 处理加入通话
func handleJoinCall(ctx context.Context, svcCtx *svc.ServiceContext, userID, roomID, callType string) {
	logx.Infof("用户加入通话: userID=%s, roomID=%s, type=%s", userID, roomID, callType)

	// 更新数据库中的用户状态为在线
	err := svcCtx.DB.Model(&call_models.CallParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("status", 1).Error // 1=在线

	if err != nil {
		logx.Errorf("更新用户通话状态失败: %v", err)
	}
}

// handleLeaveCall 处理离开通话
func handleLeaveCall(ctx context.Context, svcCtx *svc.ServiceContext, userID, roomID string) {
	logx.Infof("用户离开通话: userID=%s, roomID=%s", userID, roomID)

	// 更新数据库中的用户状态为离线
	err := svcCtx.DB.Model(&call_models.CallParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("status", 2).Error // 2=离线

	if err != nil {
		logx.Errorf("更新用户通话状态失败: %v", err)
	}
}

// handleCallHeartbeat 处理通话心跳
func handleCallHeartbeat(ctx context.Context, svcCtx *svc.ServiceContext, userID, roomID string) {
	// 更新Redis中的心跳时间戳
	heartbeatKey := fmt.Sprintf("call:heartbeat:%s:%s", roomID, userID)
	timestamp := time.Now().Unix()
	err := svcCtx.Redis.Set(heartbeatKey, timestamp, 60*time.Second).Err()
	if err != nil {
		logx.Errorf("更新通话心跳失败: %v", err)
	}

	// 确保用户状态为在线
	err = svcCtx.DB.Model(&call_models.CallParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("status", 1).Error // 1=在线

	if err != nil {
		logx.Errorf("更新用户通话状态失败: %v", err)
	}
}
