package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/livekit/protocol/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 接听通话并获取令牌
func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetCallTokenReq) (resp *types.GetCallTokenRes, err error) {
	// 1. 校验会话和参与者身份
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	isParticipant := false
	for _, pid := range session.ParticipantIds {
		if pid == req.UserID {
			isParticipant = true
			break
		}
	}

	// 核心逻辑：如果是群聊 (Type 2)，允许群内其他成员主动加入，
	// 即便他们不在初始参与者名单（ParticipantIds）中。
	if !isParticipant && session.CallType == 2 {
		isParticipant = true
	}

	if !isParticipant {
		return nil, errors.New("无权进入该通话")
	}

	// 2. 更新状态为已接听
	_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: 2, // 2-已接听
	})
	if err != nil {
		return nil, err
	}

	// 3. 生成令牌
	token, err := l.generateToken(req.UserID, req.RoomID)
	if err != nil {
		return nil, err
	}

	// 4. [核心修复] 发送信令告知所有参与者有人接听 (同步多端状态并更新其他人的成员列表)
	for _, pid := range session.ParticipantIds {
		go l.sendAcceptSignal(req.UserID, pid, req.RoomID, session.ConversationId)
	}

	// 5. 获取全量成员列表快照 (包含呼叫中、已加入、已拒绝等)
	participants := make([]types.Participant, 0)
	rpcResp, pErr := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{
		RoomId: req.RoomID,
	})
	if pErr == nil {
		for _, p := range rpcResp.Participants {
			participants = append(participants, types.Participant{
				UserID: p.UserId,
				Status: p.Status,
			})
		}
	}

	return &types.GetCallTokenRes{
		RoomToken:    token,
		LiveKitUrl:   l.svcCtx.Config.LiveKit.Host,
		Participants: participants,
	}, nil
}

func (l *GetTokenLogic) generateToken(userID, roomID string) (string, error) {
	at := auth.NewAccessToken(l.svcCtx.Config.LiveKit.ApiKey, l.svcCtx.Config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.AddGrant(grant).SetIdentity(userID).SetValidFor(time.Hour)
	return at.ToJWT()
}

func (l *GetTokenLogic) sendAcceptSignal(acceptorID, callerID, roomID, convID string) {
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
		wsCommandConst.CALL,
		wsTypeConst.CallReceive,
		acceptorID,
		callerID,
		map[string]interface{}{
			"type":   call_models.SignalAccept,
			"userId": acceptorID,
			"roomId": roomID,
		},
		convID,
	)
}
