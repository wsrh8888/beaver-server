package logic

import (
	"context"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetParticipantsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取房间当前成员列表
func NewGetParticipantsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParticipantsLogic {
	return &GetParticipantsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetParticipantsLogic) GetParticipants(req *types.GetCallParticipantsReq) (resp *types.GetCallParticipantsRes, err error) {
	rpcResp, err := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{
		RoomId: req.RoomID,
	})
	if err != nil {
		return nil, err
	}

	participants := make([]types.Participant, 0)
	for _, p := range rpcResp.Participants {
		// 统一过滤逻辑：只返回当前在该房间内（或受邀）的人
		if p.Status != 1 && p.Status != 2 {
			continue
		}

		status := "calling"
		if p.Status == 2 {
			status = "joined"
		}

		participants = append(participants, types.Participant{
			UserID: p.UserId,
			Status: status,
		})
	}

	return &types.GetCallParticipantsRes{
		Participants: participants,
	}, nil
}
