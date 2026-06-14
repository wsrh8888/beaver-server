package logic

import (
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"context"
	"encoding/json"

	"github.com/livekit/protocol/livekit"
	"github.com/zeromicro/go-zero/core/logx"
)


type LiveKitWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// LiveKit 服务器回调 (需在网关配置白名单)
func NewLiveKitWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveKitWebhookLogic {
	return &LiveKitWebhookLogic{
		ctx:    ctx,
		logger: logger.New("livekit_webhook"),
		svcCtx: svcCtx,
	}
}

func (l *LiveKitWebhookLogic) LiveKitWebhook(req *types.LiveKitWebhookReq) (resp *types.LiveKitWebhookRes, err error) {
	var event livekit.WebhookEvent
	if err := json.Unmarshal(req.Body, &event); err != nil {
		logx.WithContext(l.ctx).Errorf("解析 Webhook 事件失败: %v", err)
		return nil, err
	}

	roomID := event.Room.Name
	l.logger.Info(model.LogMsg{
		Text: "收到LiveKit回调",
		Data: map[string]interface{}{
			"event":  event.Event,
			"roomId": roomID,
		},
	})

	switch event.Event {
	case "participant_joined":
		if event.Participant != nil {
			_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
				RoomId: event.Room.Name,
				UserId: event.Participant.Identity,
				Status: int32(call_models.ParticipantStatusJoined),
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("更新参与者状态(Joined)失败: %v", err)
			}
		}
	case "participant_left":
		if event.Participant != nil {
			// 1. 更新当前人状态为已挂断/离开
			_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
				RoomId: event.Room.Name,
				UserId: event.Participant.Identity,
				Status: int32(call_models.ParticipantStatusLeft),
			})
			if err != nil {
				logx.WithContext(l.ctx).Errorf("更新参与者状态(Left)失败: %v", err)
			}

			// 2. 检查房间是否已经没有“加入中”的活跃成员了 (掉线或主动退出)
			pResp, pErr := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{RoomId: roomID})
			if pErr == nil {
				activeCount := 0
				for _, p := range pResp.Participants {
					if p.Status == int32(call_models.ParticipantStatusJoined) {
						activeCount++
					}
				}
				// 3. 如果活跃人数归零，说明大家都离开或掉线了，提前自动结束通话
				if activeCount == 0 {
					l.logger.Info(model.LogMsg{
						Text: "通话房间无活跃成员",
						Data: map[string]interface{}{"roomId": roomID},
					})
					// 标记通话结束
					_, _ = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
						RoomId: roomID,
						Status: int32(call_models.SessionStatusEnded),
					})
				}
			}
		}
	case "room_finished":
		var duration int32
		if event.Room != nil {
			duration = int32(event.CreatedAt - event.Room.CreationTime)
		}

		_, _ = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
			RoomId:   roomID,
			Duration: duration,
			Status:   int32(call_models.SessionStatusEnded),
		})
	}

	return &types.LiveKitWebhookRes{}, nil
}
