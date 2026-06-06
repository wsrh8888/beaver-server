// logic/sendmsglogic.go

package logic

import (
	"context"
	"errors"

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

// BuildMsgToRpc 将 API 的 Msg 类型转换为 RPC 的 Msg 类型（支持递归）
func (l *SendMsgLogic) BuildMsgToRpc(apiMsg *types.Msg) *chat_rpc.Msg {
	if apiMsg == nil {
		return nil
	}

	rpcMsg := &chat_rpc.Msg{
		Type: apiMsg.Type,
	}

	msgType := ctype.MsgType(apiMsg.Type)
	switch msgType {
	case ctype.TextMsgType:
		if apiMsg.TextMsg != nil {
			rpcMsg.TextMsg = &chat_rpc.TextMsg{Content: apiMsg.TextMsg.Content}
		}
	case ctype.ImageMsgType:
		if apiMsg.ImageMsg != nil {
			rpcMsg.ImageMsg = &chat_rpc.ImageMsg{
				FileKey: apiMsg.ImageMsg.FileKey,
				Width:   int32(apiMsg.ImageMsg.Width),
				Height:  int32(apiMsg.ImageMsg.Height),
				Size:    apiMsg.ImageMsg.Size,
			}
		}
	case ctype.VideoMsgType:
		if apiMsg.VideoMsg != nil {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{
				FileKey:      apiMsg.VideoMsg.FileKey,
				Width:        int32(apiMsg.VideoMsg.Width),
				Height:       int32(apiMsg.VideoMsg.Height),
				Duration:     int32(apiMsg.VideoMsg.Duration),
				ThumbnailKey: apiMsg.VideoMsg.ThumbnailKey,
				Size:         apiMsg.VideoMsg.Size,
			}
		}
	case ctype.FileMsgType:
		if apiMsg.FileMsg != nil {
			rpcMsg.FileMsg = &chat_rpc.FileMsg{
				FileKey:   apiMsg.FileMsg.FileKey,
				FileName:  apiMsg.FileMsg.FileName,
				Size:      apiMsg.FileMsg.Size,
				MimeType:  apiMsg.FileMsg.MimeType,
				Extension: apiMsg.FileMsg.Extension,
				OpenMode:  int32(apiMsg.FileMsg.OpenMode),
			}
		}
	case ctype.VoiceMsgType:
		if apiMsg.VoiceMsg != nil {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{
				FileKey:  apiMsg.VoiceMsg.FileKey,
				Duration: int32(apiMsg.VoiceMsg.Duration),
				Size:     apiMsg.VoiceMsg.Size,
			}
		}
	case ctype.EmojiMsgType:
		if apiMsg.EmojiMsg != nil {
			rpcMsg.EmojiMsg = &chat_rpc.EmojiMsg{
				FileKey:   apiMsg.EmojiMsg.FileKey,
				EmojiId:   apiMsg.EmojiMsg.EmojiID,
				PackageId: apiMsg.EmojiMsg.PackageID,
			}
		}
	case ctype.NotificationMsgType:
		if apiMsg.NotificationMsg != nil {
			rpcMsg.NotificationMsg = &chat_rpc.NotificationMsg{
				Type:   int32(apiMsg.NotificationMsg.Type),
				Actors: apiMsg.NotificationMsg.Actors,
			}
		}
	case ctype.AudioFileMsgType:
		if apiMsg.AudioFileMsg != nil {
			rpcMsg.AudioFileMsg = &chat_rpc.AudioFileMsg{
				FileKey:  apiMsg.AudioFileMsg.FileKey,
				FileName: apiMsg.AudioFileMsg.FileName,
				Duration: int32(apiMsg.AudioFileMsg.Duration),
				Size:     apiMsg.AudioFileMsg.Size,
			}
		}
	case ctype.CallMsgType:
		if apiMsg.CallMsg != nil {
			rpcMsg.CallMsg = &chat_rpc.CallMsg{
				RoomId:   apiMsg.CallMsg.RoomId,
				CallType: int32(apiMsg.CallMsg.CallType),
				Status:   int32(apiMsg.CallMsg.Status),
				Duration: apiMsg.CallMsg.Duration,
			}
		}
	case ctype.WithdrawMsgType:
		if apiMsg.WithdrawMsg != nil {
			rpcMsg.WithdrawMsg = &chat_rpc.WithdrawMsg{
				OriginMsgId: apiMsg.WithdrawMsg.OriginMsgId,
				// 递归转换快照
				OriginMsg: l.BuildMsgToRpc(apiMsg.WithdrawMsg.OriginMsg),
			}
		}
	case ctype.ReplyMsgType:
		if apiMsg.ReplyMsg != nil {
			rpcMsg.ReplyMsg = &chat_rpc.ReplyMsg{
				OriginMsgId: apiMsg.ReplyMsg.OriginMsgId,
				// 递归转换快照
				OriginMsg: l.BuildMsgToRpc(apiMsg.ReplyMsg.OriginMsg),
				// 递归转换回复的主体消息
				ReplyMsg: l.BuildMsgToRpc(apiMsg.ReplyMsg.ReplyMsg),
			}
		}
	case ctype.ForwardMsgType:
		if apiMsg.ForwardMsg != nil {
			rpcMsg.ForwardMsg = &chat_rpc.ForwardMsg{
				Title:    apiMsg.ForwardMsg.Title,
				RecordId: apiMsg.ForwardMsg.RecordID,
				Count:    int32(apiMsg.ForwardMsg.Count),
			}
		}
	case ctype.CloudDocMsgType:
		if apiMsg.CloudDocMsg != nil {
			rpcMsg.CloudDocMsg = &chat_rpc.CloudDocMsg{
				DocId: apiMsg.CloudDocMsg.DocID,
				DocType:  int32(apiMsg.CloudDocMsg.DocType),
				Title:    apiMsg.CloudDocMsg.Title,
				OwnerId:  apiMsg.CloudDocMsg.OwnerID,
				Perm:     int32(apiMsg.CloudDocMsg.Perm),
				CoverUrl: apiMsg.CloudDocMsg.CoverURL,
				Revision: apiMsg.CloudDocMsg.Revision,
			}
		}
	}
	return rpcMsg
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgReq) (*types.SendMsgRes, error) {
	// 构建 RPC 请求的消息内容（使用递归转换逻辑）
	rpcMsg := l.BuildMsgToRpc(&req.Msg)

	rpcReq := &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		MessageId:      req.MessageID,
		ConversationId: req.ConversationID,
		Msg:            rpcMsg,
	}

	// 记录请求（仅调试用）
	l.Logger.Infof("Sending message via RPC: userId=%s, conversationId=%s, type=%d", req.UserID, req.ConversationID, req.Msg.Type)

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("failed to send message via RPC: %v", err)
		return nil, errors.New("failed to send message")
	}

	// 构建 API 响应
	resp := &types.SendMsgRes{
		Id:             uint(rpcResp.Id),
		MessageID:      rpcResp.MessageId,
		ConversationID: rpcResp.ConversationId,
		Msg:            req.Msg, // 返回原始发送的消息对象
		Sender: types.Sender{
			UserID:   rpcResp.Sender.UserId,
			Avatar:   rpcResp.Sender.Avatar,
			NickName: rpcResp.Sender.NickName,
			UserType: int8(rpcResp.Sender.UserType),
		},
		CreatedAt:  rpcResp.CreatedAt,
		MsgPreview: rpcResp.MsgPreview,
		Seq:        rpcResp.Seq,
	}

	return resp, nil
}
