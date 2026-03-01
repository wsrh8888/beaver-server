package logic

import (
	"context"
	"errors"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

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

	// 4. 发送信令告知其他人有人挂断/离开 (纯信令，不入库)
	for _, pid := range session.ParticipantIds {
		if pid != req.UserID {
			go l.sendHangupSignal(req.UserID, pid, req.RoomID, session.ConversationId)
		}
	}

	return &types.HangupCallRes{}, nil
}

func (l *HangupLogic) sendHangupSignal(hanguperID, targetID, roomID, convID string) {
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
		convID,
	)
}
