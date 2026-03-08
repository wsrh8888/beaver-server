package logic

import (
	"context"
	"time"

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/chat"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type FinalizeSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFinalizeSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinalizeSessionLogic {
	return &FinalizeSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 结束通话记录
func (l *FinalizeSessionLogic) FinalizeSession(in *call_rpc.FinalizeSessionReq) (*call_rpc.FinalizeSessionRes, error) {
	// 1. 获取会话原始信息 (为了拿到锚点 MessageID 和会话ID)
	var session call_models.CallSession
	if err := l.svcCtx.DB.Where("room_id = ?", in.RoomId).First(&session).Error; err != nil {
		return nil, err
	}

	// 1.1 防重校验：如果状态已经是结束状态(3,4,5)，直接返回成功，不再发补丁
	if session.Status == call_models.SessionStatusEnded ||
		session.Status == call_models.SessionStatusMissed ||
		session.Status == call_models.SessionStatusRejected {
		return &call_rpc.FinalizeSessionRes{Success: true}, nil
	}

	// 2. 更新数据库状态
	now := time.Now()
	err := l.svcCtx.DB.Model(&session).Updates(map[string]interface{}{
		"status":   int8(in.Status),
		"end_time": &now,
		"duration": in.Duration,
	}).Error
	if err != nil {
		return nil, err
	}

	// 3. 发送聊天记录“状态补丁”消息 (基于 TargetMsgID 追加法)
	// 使用 RoomID 作为锚点消息的 ID (因为我们在发起时将 MessageID 设置为了 RoomID)
	_, _ = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat.SendMsgReq{
		UserId:         session.CallerID,
		ConversationId: session.ConversationID,
		MessageId:      uuid.New().String(),
		Msg: &chat.Msg{
			Type:        9,              // CallMsg
			TargetMsgId: session.RoomID, // 指向发起时的那条消息 (RoomID == Initial MessageID)
			CallMsg: &chat.CallMsg{
				RoomId:   session.RoomID,
				CallType: int32(session.CallType),
				Status:   2, // 2-已结束
				Duration: int64(in.Duration),
			},
		},
	})

	return &call_rpc.FinalizeSessionRes{Success: true}, nil
}
