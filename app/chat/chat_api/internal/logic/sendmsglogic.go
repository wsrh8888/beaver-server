// logic/sendmsglogic.go

package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgReq) (*types.SendMsgRes, error) {
	// 构建RPC请求
	rpcReq := &chat_rpc.SendMsgReq{
		UserID:         req.UserID,
		ConversationId: req.ConversationID,
		Msg: &chat_rpc.Msg{
			Type: req.Msg.Type,
		},
	}
	fmt.Println("re1111111111111111q:", req)
	msgType := ctype.MsgType(req.Msg.Type)
	switch msgType {
	case ctype.TextMsgType:
		rpcReq.Msg.TextMsg = &chat_rpc.TextMsg{Content: req.Msg.TextMsg.Content}
	case ctype.ImageMsgType:
		rpcReq.Msg.ImageMsg = &chat_rpc.ImageMsg{Name: req.Msg.ImageMsg.Name, FileId: req.Msg.ImageMsg.FileId}
	case ctype.VideoMsgType:
		rpcReq.Msg.VideoMsg = &chat_rpc.VideoMsg{Title: req.Msg.VideoMsg.Title, Src: req.Msg.VideoMsg.Src, Time: req.Msg.VideoMsg.Time}
	case ctype.FileMsgType:
		rpcReq.Msg.FileMsg = &chat_rpc.FileMsg{Title: req.Msg.FileMsg.Title, Src: req.Msg.FileMsg.Src, Size: req.Msg.FileMsg.Size, Type: req.Msg.FileMsg.Type}
	case ctype.VoiceMsgType:
		rpcReq.Msg.VoiceMsg = &chat_rpc.VoiceMsg{Src: req.Msg.VoiceMsg.Src, Time: req.Msg.VoiceMsg.Time}
	default:
		return nil, errors.New("invalid message type")

	}
	fmt.Println("rpcReq:", rpcReq)
	// 调用RPC服务
	rpcResp, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("failed to send message via RPC: %v", err)
		return nil, errors.New("failed to send message")
	}

	// 构建API响应
	resp := &types.SendMsgRes{
		MessageID:      rpcResp.MessageId,
		ConversationID: rpcResp.ConversationId,
		Msg:            req.Msg,
		Sender: types.Sender{
			UserID:   rpcResp.Sender.UserID,
			Avatar:   rpcResp.Sender.Avatar,
			Nickname: rpcResp.Sender.Nickname,
		},
		CreateAt:   rpcResp.CreateAt,
		MsgPreview: rpcResp.MsgPreview,
	}

	return resp, nil
}
