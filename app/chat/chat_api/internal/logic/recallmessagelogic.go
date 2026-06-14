package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/models/ctype"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)


type RecallMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewRecallMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallMessageLogic {
	return &RecallMessageLogic{
		ctx:    ctx,
		logger: logger.New("recall_msg"),
		svcCtx: svcCtx,
	}
}

func (l *RecallMessageLogic) RecallMessage(req *types.RecallMessageReq) (resp *types.RecallMessageRes, err error) {
	// 1. 获取原始消息
	var msg chat_models.ChatMessage
	err = l.svcCtx.DB.Where("message_id = ?", req.MessageID).First(&msg).Error
	if err != nil {
		return nil, errors.New("消息不存在")
	}

	// 2. 权限校验
	if msg.SendUserID == nil || *msg.SendUserID != req.UserID {
		return nil, errors.New("无权撤回他人消息")
	}

	// 3. 时效性校验（对标大厂：3分钟限制）
	if time.Since(time.Time(msg.CreatedAt)) > 3*time.Minute {
		return nil, errors.New("超过3分钟，无法撤回")
	}

	// 4. 不修改原始消息记录（只增不改原则）
	// 撤回通过发送一条新的 WithdrawMsg 指令消息来表达
	// 这条指令消息有自己的 Seq，会正常进入同步流，客户端通过 originMsgId 识别被撤回的消息

	// 5. 产生一条"撤回指令"到同步流
	// 这里调用 SendMsg 发送一条 WithdrawMsg 类型的消息，它会产生新的 Seq
	_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		ConversationId: msg.ConversationID,
		MessageId:      uuid.New().String(), // 指令消息本身的ID
		Msg: &chat_rpc.Msg{
			Type: uint32(ctype.WithdrawMsgType),
			WithdrawMsg: &chat_rpc.WithdrawMsg{
				OriginMsgId: req.MessageID,
				// 将被撤回的消息快照带上，供展示或审计
				OriginMsg: func() *chat_rpc.Msg {
					res := &chat_rpc.Msg{Type: uint32(msg.MsgType)}
					if msg.MsgType == ctype.TextMsgType && msg.Msg != nil && msg.Msg.TextMsg != nil {
						res.TextMsg = &chat_rpc.TextMsg{Content: msg.Msg.TextMsg.Content}
					}
					return res
				}(),
			},
		},
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("发送撤回指令失败: %v", err)
		return nil, errors.New("撤回失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "消息撤回成功",
		Data: map[string]interface{}{
			"messageId":      req.MessageID,
			"userId":         req.UserID,
			"conversationId": msg.ConversationID,
		},
	})

	return &types.RecallMessageRes{
		MessageID:  req.MessageID,
		RecallTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
