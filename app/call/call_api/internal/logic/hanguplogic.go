package logic

import (
	"context"
	"errors"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type HangupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 主动挂断或拒绝通话
func NewHangupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HangupLogic {
	return &HangupLogic{
		ctx:    ctx,
		logger: logger.New("hangup_call"),
		svcCtx: svcCtx,
	}
}

func (l *HangupLogic) Hangup(req *types.HangupCallReq) (resp *types.HangupCallRes, err error) {
	// 1. 获取会话信息
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	// 2. 根据通话类型决定状态和信令
	status := int32(5) // 默认挂断 (Left)
	signalType := call_models.SignalHangup

	if session.CallType == 2 {
		// 群聊：统一视为拒绝/退出，不结束会话
		status = 3 // 拒绝 (Rejected)
		signalType = call_models.SignalReject
	}

	// 3. 更新当前参与者状态
	_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: status,
	})

	// 4. 判断是否需要结束整个会话 (私聊或发起者主动挂断)
	if session.CallType == 1 || req.UserID == session.CallerId {
		_, _ = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
			RoomId: req.RoomID,
			Status: 3, // 3-已结束
		})
	}

	// 5. [核心修复] 发送信令告知所有人 (包括自己的其他设备同步)
	for _, pid := range session.ParticipantIds {
		go l.sendSignal(req.UserID, pid, req.RoomID, session.ConversationId, signalType)
	}

	l.logger.Info(model.LogMsg{
		Text: "通话挂断成功",
		Data: map[string]interface{}{
			"roomId": req.RoomID,
			"userId": req.UserID,
			"status": status,
		},
	})

	return &types.HangupCallRes{}, nil
}

func (l *HangupLogic) sendSignal(hanguperID, targetID, roomID, convID, signalType string) {
	payload := map[string]interface{}{
		"command":  wsCommandConst.CALL,
		"type":     wsTypeConst.CallReceive,
		"senderId": hanguperID,
		"targetId": targetID,
		"body": map[string]interface{}{
			"type":   signalType,
			"userId": hanguperID,
			"roomId": roomID,
		},
		"conversationId": convID,
	}
	l.svcCtx.RocketMQ.SendMessage(l.ctx, mqwsconst.MqTopicWs, payload)
}
