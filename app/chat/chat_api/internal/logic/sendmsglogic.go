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
		UserId:         req.UserID,
		MessageId:      req.MessageID,
		ConversationId: req.ConversationID,
		Msg: &chat_rpc.Msg{
			Type: req.Msg.Type,
		},
	}
	msgType := ctype.MsgType(req.Msg.Type)
	switch msgType {
	case ctype.TextMsgType:
		rpcReq.Msg.TextMsg = &chat_rpc.TextMsg{Content: req.Msg.TextMsg.Content}
	case ctype.ImageMsgType:
		imageMsg := &chat_rpc.ImageMsg{FileKey: req.Msg.ImageMsg.FileKey}
		// 如果存在 style，则设置 style
		if req.Msg.ImageMsg.Style != nil {
			imageMsg.Style = &chat_rpc.ImageStyle{
				Width:  int32(req.Msg.ImageMsg.Style.Width),
				Height: int32(req.Msg.ImageMsg.Style.Height),
			}
		}
		rpcReq.Msg.ImageMsg = imageMsg
	case ctype.VideoMsgType:
		videoMsg := &chat_rpc.VideoMsg{FileKey: req.Msg.VideoMsg.FileKey}
		// 如果存在 style，则设置 style
		if req.Msg.VideoMsg.Style != nil {
			videoMsg.Style = &chat_rpc.VideoStyle{
				Width:    int32(req.Msg.VideoMsg.Style.Width),
				Height:   int32(req.Msg.VideoMsg.Style.Height),
				Duration: int32(req.Msg.VideoMsg.Style.Duration),
			}
		}
		rpcReq.Msg.VideoMsg = videoMsg
	case ctype.FileMsgType:
		rpcReq.Msg.FileMsg = &chat_rpc.FileMsg{FileKey: req.Msg.FileMsg.FileKey}
	case ctype.VoiceMsgType:
		voiceMsg := &chat_rpc.VoiceMsg{FileKey: req.Msg.VoiceMsg.FileKey}
		// 如果存在 style，则设置 style
		if req.Msg.VoiceMsg.Style != nil {
			voiceMsg.Style = &chat_rpc.VoiceStyle{
				Duration: int32(req.Msg.VoiceMsg.Style.Duration),
			}
		}
		rpcReq.Msg.VoiceMsg = voiceMsg
	case ctype.EmojiMsgType:
		rpcReq.Msg.EmojiMsg = &chat_rpc.EmojiMsg{
			FileKey:   req.Msg.EmojiMsg.FileKey,
			EmojiId:   req.Msg.EmojiMsg.EmojiID,
			PackageId: req.Msg.EmojiMsg.PackageID,
		}
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
		Id:             uint(rpcResp.Id),
		ConversationID: rpcResp.ConversationId,
		Msg:            req.Msg,
		Sender: types.Sender{
			UserID:   rpcResp.Sender.UserId,
			Avatar:   rpcResp.Sender.Avatar,
			Nickname: rpcResp.Sender.Nickname,
		},
		CreateAt:   rpcResp.CreateAt,
		MsgPreview: rpcResp.MsgPreview,
		Status:     uint32(rpcResp.Status),
		Seq:        rpcResp.Seq,
	}

	return resp, nil
}
