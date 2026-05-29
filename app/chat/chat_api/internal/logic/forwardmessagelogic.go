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

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ForwardMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForwardMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForwardMessageLogic {
	return &ForwardMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForwardMessageLogic) ForwardMessage(req *types.ForwardMessageReq) (resp *types.ForwardMessageRes, err error) {
	// 1. 获取原始消息对象列表
	var originMessages []chat_models.ChatMessage
	err = l.svcCtx.DB.Where("message_id IN ?", req.MessageIDs).Order("created_at asc").Find(&originMessages).Error
	if err != nil {
		l.Logger.Errorf("获取待转发消息失败: %v", err)
		return nil, err
	}

	if len(originMessages) == 0 {
		return nil, errors.New("未找到有效的转发消息")
	}

	if req.ForwardMode == 1 {
		// --- 逐条转发模式 ---
		for _, m := range originMessages {
			// 生成全新的客户端ID（简单起见使用UUID+原始ID，大厂通常由前端传入或后端补齐）
			newMsgID := uuid.New().String()

			// 调用 RPC 发送消息 (复用发送逻辑)
			_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
				UserId:         req.UserID,
				ConversationId: req.TargetID,
				MessageId:      newMsgID,
				Msg:            l.convertModelToProtoMsg(m.Msg),
			})
			if err != nil {
				l.Logger.Errorf("逐条转发失败: %v", err)
				// 商业化项目通常会继续处理下一条，或者返回部分成功的提示
			}
		}
	} else {
		// --- 合并转发模式 ---
		recordID := uuid.New().String()

		// 2. 将消息快照存入详情表 (冷数据)
		err = l.svcCtx.DB.Create(&chat_models.ChatForward{
			RecordID: recordID,
			Content:  originMessages, // 直接赋值，由 ForwardContent.Value 接口处理序列化
		}).Error
		if err != nil {
			l.Logger.Errorf("创建转发详情失败: %v", err)
			return nil, err
		}

		// 3. 发送合并转发卡片 (热数据)
		title := "聊天记录"
		if len(originMessages) > 0 {
			title = "群聊的聊天记录"
		}

		_, err = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
			UserId:         req.UserID,
			ConversationId: req.TargetID,
			MessageId:      uuid.New().String(),
			Msg: &chat_rpc.Msg{
				Type: uint32(ctype.ForwardMsgType),
				ForwardMsg: &chat_rpc.ForwardMsg{
					Title:    title,
					RecordId: recordID,
					Count:    int32(len(originMessages)),
				},
			},
		})
		if err != nil {
			l.Logger.Errorf("发送合并转发卡片失败: %v", err)
			return nil, err
		}
	}

	return &types.ForwardMessageRes{
		ForwardTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// 辅助方法：将 DB Model 的 Msg 转换为 RPC 的 Msg (需要适配具体字段)
func (l *ForwardMessageLogic) convertModelToProtoMsg(m *ctype.Msg) *chat_rpc.Msg {
	if m == nil {
		return nil
	}

	rpcMsg := &chat_rpc.Msg{
		Type: uint32(m.Type),
	}

	switch m.Type {
	case ctype.TextMsgType:
		if m.TextMsg != nil {
			rpcMsg.TextMsg = &chat_rpc.TextMsg{Content: m.TextMsg.Content}
		}
	case ctype.ImageMsgType:
		if m.ImageMsg != nil {
			rpcMsg.ImageMsg = &chat_rpc.ImageMsg{
				FileKey: m.ImageMsg.FileKey,
				Width:   int32(m.ImageMsg.Width),
				Height:  int32(m.ImageMsg.Height),
				Size:    m.ImageMsg.Size,
			}
		}
	case ctype.VideoMsgType:
		if m.VideoMsg != nil {
			rpcMsg.VideoMsg = &chat_rpc.VideoMsg{
				FileKey:      m.VideoMsg.FileKey,
				Width:        int32(m.VideoMsg.Width),
				Height:       int32(m.VideoMsg.Height),
				Duration:     int32(m.VideoMsg.Duration),
				ThumbnailKey: m.VideoMsg.ThumbnailKey,
				Size:         m.VideoMsg.Size,
			}
		}
	case ctype.FileMsgType:
		if m.FileMsg != nil {
			rpcMsg.FileMsg = &chat_rpc.FileMsg{
				FileKey:   m.FileMsg.FileKey,
				FileName:  m.FileMsg.FileName,
				Size:      m.FileMsg.Size,
				MimeType:  m.FileMsg.MimeType,
				Extension: m.FileMsg.Extension,
				OpenMode:  int32(m.FileMsg.OpenMode),
			}
		}
	case ctype.VoiceMsgType:
		if m.VoiceMsg != nil {
			rpcMsg.VoiceMsg = &chat_rpc.VoiceMsg{
				FileKey:  m.VoiceMsg.FileKey,
				Duration: int32(m.VoiceMsg.Duration),
				Size:     m.VoiceMsg.Size,
			}
		}
	case ctype.EmojiMsgType:
		if m.EmojiMsg != nil {
			rpcMsg.EmojiMsg = &chat_rpc.EmojiMsg{
				FileKey:   m.EmojiMsg.FileKey,
				EmojiId:   m.EmojiMsg.EmojiID,
				PackageId: m.EmojiMsg.PackageID,
				Width:     m.EmojiMsg.Width,
				Height:    m.EmojiMsg.Height,
			}
		}
	case ctype.NotificationMsgType:
		if m.NotificationMsg != nil {
			rpcMsg.NotificationMsg = &chat_rpc.NotificationMsg{
				Type:   int32(m.NotificationMsg.Type),
				Actors: m.NotificationMsg.Actors,
			}
		}
	case ctype.AudioFileMsgType:
		if m.AudioFileMsg != nil {
			rpcMsg.AudioFileMsg = &chat_rpc.AudioFileMsg{
				FileKey:  m.AudioFileMsg.FileKey,
				FileName: m.AudioFileMsg.FileName,
				Duration: int32(m.AudioFileMsg.Duration),
				Size:     m.AudioFileMsg.Size,
			}
		}
	case ctype.CallMsgType:
		if m.CallMsg != nil {
			rpcMsg.CallMsg = &chat_rpc.CallMsg{
				RoomId:   m.CallMsg.RoomID,
				CallType: int32(m.CallMsg.CallType),
				Status:   int32(m.CallMsg.Status),
				Duration: m.CallMsg.Duration,
			}
		}
	case ctype.WithdrawMsgType:
		if m.WithdrawMsg != nil {
			rpcMsg.WithdrawMsg = &chat_rpc.WithdrawMsg{
				OriginMsgId: m.WithdrawMsg.OriginMsgID,
				OriginMsg:   l.convertModelToProtoMsg(m.WithdrawMsg.OriginMsg),
			}
		}
	case ctype.ReplyMsgType:
		if m.ReplyMsg != nil {
			rpcMsg.ReplyMsg = &chat_rpc.ReplyMsg{
				OriginMsgId: m.ReplyMsg.OriginMsgID,
				OriginMsg:   l.convertModelToProtoMsg(m.ReplyMsg.OriginMsg),
				ReplyMsg:    l.convertModelToProtoMsg(m.ReplyMsg.ReplyMsg),
			}
		}
	case ctype.ForwardMsgType:
		if m.ForwardMsg != nil {
			rpcMsg.ForwardMsg = &chat_rpc.ForwardMsg{
				Title:    m.ForwardMsg.Title,
				RecordId: m.ForwardMsg.RecordID,
				Count:    int32(m.ForwardMsg.Count),
			}
		}
	case ctype.MarkdownMsgType:
		if m.MarkdownMsg != nil {
			rpcMsg.MarkdownMsg = &chat_rpc.MarkdownMsg{
				Content: m.MarkdownMsg.Content,
				Title:   m.MarkdownMsg.Title,
			}
		}
	case ctype.LinkMsgType:
		if m.LinkMsg != nil {
			rpcMsg.LinkMsg = &chat_rpc.LinkMsg{
				Url:      m.LinkMsg.URL,
				Title:    m.LinkMsg.Title,
				Desc:     m.LinkMsg.Desc,
				ImageUrl: m.LinkMsg.ImageURL,
			}
		}
	case ctype.CloudDocMsgType:
		if m.CloudDocMsg != nil {
			rpcMsg.CloudDocMsg = &chat_rpc.CloudDocMsg{
				DocId: m.CloudDocMsg.DocID,
				DocType:  int32(m.CloudDocMsg.DocType),
				Title:    m.CloudDocMsg.Title,
				OwnerId:  m.CloudDocMsg.OwnerID,
				Perm:     int32(m.CloudDocMsg.Perm),
				CoverUrl: m.CloudDocMsg.CoverURL,
				Revision: m.CloudDocMsg.Revision,
			}
		}
	}
	return rpcMsg
}
