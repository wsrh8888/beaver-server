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
		// 设置可选字段（打平后的结构）
		if req.Msg.ImageMsg.Width > 0 {
			imageMsg.Width = int32(req.Msg.ImageMsg.Width)
		}
		if req.Msg.ImageMsg.Height > 0 {
			imageMsg.Height = int32(req.Msg.ImageMsg.Height)
		}
		if req.Msg.ImageMsg.Size > 0 {
			imageMsg.Size = req.Msg.ImageMsg.Size
		}
		rpcReq.Msg.ImageMsg = imageMsg
	case ctype.VideoMsgType:
		videoMsg := &chat_rpc.VideoMsg{FileKey: req.Msg.VideoMsg.FileKey}
		// 设置可选字段（打平后的结构）
		if req.Msg.VideoMsg.Width > 0 {
			videoMsg.Width = int32(req.Msg.VideoMsg.Width)
		}
		if req.Msg.VideoMsg.Height > 0 {
			videoMsg.Height = int32(req.Msg.VideoMsg.Height)
		}
		if req.Msg.VideoMsg.Duration > 0 {
			videoMsg.Duration = int32(req.Msg.VideoMsg.Duration)
		}
		if req.Msg.VideoMsg.ThumbnailKey != "" {
			videoMsg.ThumbnailKey = req.Msg.VideoMsg.ThumbnailKey
		}
		if req.Msg.VideoMsg.Size > 0 {
			videoMsg.Size = req.Msg.VideoMsg.Size
		}
		rpcReq.Msg.VideoMsg = videoMsg
	case ctype.FileMsgType:
		fileMsg := &chat_rpc.FileMsg{FileKey: req.Msg.FileMsg.FileKey}
		// 设置可选字段
		if req.Msg.FileMsg.FileName != "" {
			fileMsg.FileName = req.Msg.FileMsg.FileName
		}
		if req.Msg.FileMsg.Size > 0 {
			fileMsg.Size = req.Msg.FileMsg.Size
		}
		if req.Msg.FileMsg.MimeType != "" {
			fileMsg.MimeType = req.Msg.FileMsg.MimeType
		}
		rpcReq.Msg.FileMsg = fileMsg
	case ctype.VoiceMsgType:
		voiceMsg := &chat_rpc.VoiceMsg{FileKey: req.Msg.VoiceMsg.FileKey}
		// 设置可选字段（打平后的结构）
		if req.Msg.VoiceMsg.Duration > 0 {
			voiceMsg.Duration = int32(req.Msg.VoiceMsg.Duration)
		}
		if req.Msg.VoiceMsg.Size > 0 {
			voiceMsg.Size = req.Msg.VoiceMsg.Size
		}
		rpcReq.Msg.VoiceMsg = voiceMsg
	case ctype.EmojiMsgType:
		rpcReq.Msg.EmojiMsg = &chat_rpc.EmojiMsg{
			FileKey:   req.Msg.EmojiMsg.FileKey,
			EmojiId:   req.Msg.EmojiMsg.EmojiID,
			PackageId: req.Msg.EmojiMsg.PackageID,
		}
	case ctype.NotificationMsgType:
		rpcReq.Msg.NotificationMsg = &chat_rpc.NotificationMsg{
			Type:   int32(req.Msg.NotificationMsg.Type),
			Actors: req.Msg.NotificationMsg.Actors,
		}
	case ctype.AudioFileMsgType:
		audioFileMsg := &chat_rpc.AudioFileMsg{FileKey: req.Msg.AudioFileMsg.FileKey}
		// 设置可选字段（打平后的结构）
		if req.Msg.AudioFileMsg.FileName != "" {
			audioFileMsg.FileName = req.Msg.AudioFileMsg.FileName
		}
		if req.Msg.AudioFileMsg.Duration > 0 {
			audioFileMsg.Duration = int32(req.Msg.AudioFileMsg.Duration)
		}
		if req.Msg.AudioFileMsg.Size > 0 {
			audioFileMsg.Size = req.Msg.AudioFileMsg.Size
		}
		rpcReq.Msg.AudioFileMsg = audioFileMsg
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
