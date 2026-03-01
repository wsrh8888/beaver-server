package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type HangupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 主动挂断或拒绝通话
func NewHangupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HangupLogic {
	return &HangupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HangupLogic) Hangup(req *types.HangupCallReq) (resp *types.HangupCallRes, err error) {
	// 1. 获取会话信息
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	// 2. 更新当前参与者状态为已挂断
	_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: 5, // 5-已挂断
	})

	// 3. 判断是否需要结束整个会话
	shouldFinalize := false
	if session.CallType == 1 {
		// 私聊：任意一方挂断，整个通话结束
		shouldFinalize = true
	} else {
		// 群聊：检查是否还有其他人在房间里 (joined 状态)
		participants, pErr := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{RoomId: req.RoomID})
		if pErr == nil {
			activeCount := 0
			for _, p := range participants.Participants {
				if p.Status == 2 { // 2-已接听/加入中
					activeCount++
				}
			}
			// 如果房间里已经没有其他人了（刚才那个人已经退出了，所以这里 activeCount 为 0 则代表全空）
			if activeCount == 0 {
				shouldFinalize = true
			}
		}
	}

	if shouldFinalize {
		_, _ = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
			RoomId: req.RoomID,
			Status: 3, // 3-已结束
		})
	}

	// 4. 发送信令告知其他人有人挂断/离开
	for _, pid := range session.ParticipantIds {
		if pid != req.UserID {
			go l.sendHangupSignal(req.UserID, pid, req.RoomID)
		}
	}

	return &types.HangupCallRes{}, nil
}

func (l *HangupLogic) sendHangupSignal(hanguperID, targetID, roomID string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type":   "RTC_HANGUP",
		"user":   hanguperID,
		"roomId": roomID,
	})

	_, err := l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         hanguperID,
		ConversationId: l.getConversationID(hanguperID, targetID),
		MessageId:      uuid.New().String(), // 注入唯一 ID
		Msg: &chat_rpc.Msg{
			Type: 7, // 7:通知消息/信令
			NotificationMsg: &chat_rpc.NotificationMsg{
				Type:   102, // RTC_HANGUP
				Actors: []string{hanguperID},
			},
			TextMsg: &chat_rpc.TextMsg{
				Content: string(payload),
			},
		},
	})
	if err != nil {
		logx.Errorf("发送 RTC_HANGUP 信令失败: %v", err)
	}

	// 2. 直接通过 WebSocket 发送 RTC 信令
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
		wsCommandConst.CALL,
		wsTypeConst.CallReceive,
		hanguperID,
		targetID,
		map[string]interface{}{
			"type":   call_models.SignalHangup,
			"user":   hanguperID,
			"roomId": roomID,
		},
		l.getConversationID(hanguperID, targetID),
	)
}

func (l *HangupLogic) getConversationID(callerID, targetID string) string {
	// 如果 targetID 是群 ID
	if len(targetID) > 6 && targetID[:6] == "group_" {
		return targetID
	}
	// 私聊会话 ID 拼装
	if callerID < targetID {
		return callerID + ":" + targetID
	}
	return targetID + ":" + callerID
}
