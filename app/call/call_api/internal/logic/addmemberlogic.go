package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/livekit/protocol/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群聊中途加入通话
func NewAddMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberLogic {
	return &AddMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMemberLogic) AddMember(req *types.AddCallMemberReq) (resp *types.AddCallMemberRes, err error) {
	// 1. 校验会话
	_, err = l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	// 2. 更新/加入状态为已接听 (joined)
	_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: 2, // 2-已接听/joined
	})
	if err != nil {
		return nil, err
	}

	// 3. 生成令牌
	at := auth.NewAccessToken(l.svcCtx.Config.LiveKit.ApiKey, l.svcCtx.Config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     req.RoomID,
	}
	at.AddGrant(grant).SetIdentity(req.UserID).SetValidFor(time.Hour)
	token, err := at.ToJWT()
	if err != nil {
		return nil, err
	}

	return &types.AddCallMemberRes{
		RoomToken:  token,
		LiveKitUrl: l.svcCtx.Config.LiveKit.Host,
	}, nil
}
