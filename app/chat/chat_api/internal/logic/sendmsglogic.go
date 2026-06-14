// logic/sendmsglogic.go

package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/models/ctype"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		logger: logger.New("send_msg"),
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
				FileUrl: apiMsg.ImageMsg.FileUrl,
				Width:   int32(apiMsg.ImageMsg.Width),
				Height:  int32(apiMsg.ImageMsg.Height),
				Size:    apiMsg.ImageMsg.Size,
			}
		}
	case ctype.VideoMsgType:
		if apiMsg.VideoMsg != nil {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{
				FileUrl:       apiMsg.VideoMsg.FileUrl,
				Width:        int32(apiMsg.VideoMsg.Width),
				Height:       int32(apiMsg.VideoMsg.Height),
				Duration:     int32(apiMsg.VideoMsg.Duration),
				ThumbnailUrl: apiMsg.VideoMsg.ThumbnailUrl,
				Size:         apiMsg.VideoMsg.Size,
			}
		}
	case ctype.FileMsgType:
		if apiMsg.FileMsg != nil {
			rpcMsg.FileMsg = &chat_rpc.FileMsg{
				FileUrl:   apiMsg.FileMsg.FileUrl,
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
				FileUrl:  apiMsg.VoiceMsg.FileUrl,
				Duration: int32(apiMsg.VoiceMsg.Duration),
				Size:     apiMsg.VoiceMsg.Size,
			}
		}
	case ctype.EmojiMsgType:
		if apiMsg.EmojiMsg != nil {
			rpcMsg.EmojiMsg = &chat_rpc.EmojiMsg{
				FileUrl:   apiMsg.EmojiMsg.FileUrl,
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
				FileUrl:  apiMsg.AudioFileMsg.FileUrl,
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

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, rpcReq)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("failed to send message via RPC: %v", err)
		l.logger.Error(model.LogMsg{
			Text: "发送消息失败",
			Data: map[string]interface{}{
				"userId":         req.UserID,
				"conversationId": req.ConversationID,
				"messageType":    req.Msg.Type,
			},
		})
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

	l.logger.Info(model.LogMsg{
		Text: "发送消息成功",
		Data: map[string]interface{}{
			"userId":         req.UserID,
			"conversationId": req.ConversationID,
			"messageId":      rpcResp.MessageId,
			"messageType":    req.Msg.Type,
		},
	})
	return resp, nil
}
