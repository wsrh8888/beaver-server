package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/chat/chat_utils"
	"beaver/app/friend/friend_models"
	"beaver/common/models/ctype"
	"beaver/utils/conversation"

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

func (l *SendMsgLogic) SendMsg(in *chat_rpc.SendMsgReq) (*chat_rpc.SendMsgRes, error) {
	if conversation.GetConversationType(in.ConversationId) == 1 {
		if !strings.Contains(in.ConversationId, in.UserId) {
			logx.Errorf("用户id不匹配，用户id：%s，会话id：%s", in.UserId, in.ConversationId)
			return nil, errors.New("异常操作")
		}
		var friend friend_models.FriendModel
		userIds := conversation.ParseConversation(in.ConversationId)
		if !friend.IsFriend(l.svcCtx.DB, userIds[0], userIds[1]) {
			return nil, errors.New("不是好友关系")
		}
	}

	msgType := ctype.MsgType(in.Msg.Type)
	var msg ctype.Msg
	switch msgType {
	case ctype.TextMsgType:
		msg = ctype.Msg{
			Type: ctype.TextMsgType,
			TextMsg: &ctype.TextMsg{
				Content: in.Msg.TextMsg.Content,
			},
		}
	case ctype.ImageMsgType:
		msg = ctype.Msg{
			Type: ctype.ImageMsgType,
			ImageMsg: &ctype.ImageMsg{
				FileKey: in.Msg.ImageMsg.FileKey,
				Width:   int(in.Msg.ImageMsg.Width),
				Height:  int(in.Msg.ImageMsg.Height),
			},
		}
	case ctype.VideoMsgType:
		msg = ctype.Msg{
			Type: ctype.VideoMsgType,
			VideoMsg: &ctype.VideoMsg{
				FileKey:  in.Msg.VideoMsg.FileKey,
				Width:    int(in.Msg.VideoMsg.Width),
				Height:   int(in.Msg.VideoMsg.Height),
				Duration: int(in.Msg.VideoMsg.Duration),
			},
		}
	case ctype.FileMsgType:
		msg = ctype.Msg{
			Type: ctype.FileMsgType,
			FileMsg: &ctype.FileMsg{
				FileKey: in.Msg.FileMsg.FileKey,
			},
		}
	case ctype.VoiceMsgType:
		msg = ctype.Msg{
			Type: ctype.VoiceMsgType,
			VoiceMsg: &ctype.VoiceMsg{
				FileKey:  in.Msg.VoiceMsg.FileKey,
				Duration: int(in.Msg.VoiceMsg.Duration),
			},
		}
	case ctype.EmojiMsgType:
		msg = ctype.Msg{
			Type: ctype.EmojiMsgType,
			EmojiMsg: &ctype.EmojiMsg{
				FileKey:   in.Msg.EmojiMsg.FileKey,
				EmojiID:   in.Msg.EmojiMsg.EmojiId,
				PackageID: in.Msg.EmojiMsg.PackageId,
			},
		}
	default:
		return nil, fmt.Errorf("未识别到该类型: %d", msgType)
	}

	chatModel := chat_models.ChatMessage{
		SendUserID:     in.UserId,
		MessageID:      in.MessageId,
		ConversationID: in.ConversationId,
		MsgType:        msgType,
		Msg:            &msg,
	}

	// 1. 创建消息记录
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	err := l.svcCtx.DB.Create(&chatModel).Error
	if err != nil {
		return nil, err
	}

	// 设置序列号为数据库自增ID（全局递增）
	chatModel.Seq = int64(chatModel.Id)
	// 更新数据库中的seq字段
	err = l.svcCtx.DB.Model(&chatModel).Update("seq", chatModel.Seq).Error
	if err != nil {
		return nil, err
	}

	// 2. 更新会话级别的信息
	err = chat_utils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversation.GetConversationType(in.ConversationId), chatModel.Seq)
	if err != nil {
		return nil, err
	}

	// 3. 更新用户会话关系
	err = chat_utils.UpdateUserConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.UserId, in.ConversationId, chatModel.MsgPreview, false)
	if err != nil {
		return nil, err
	}

	// 注意：同步游标由客户端主动更新，服务端不主动更新
	// 客户端可以通过以下方式获取更新：
	// 1. 定期轮询同步接口
	// 2. 通过WebSocket推送通知
	// 3. 客户端主动调用 UpdateSyncCursor API

	convertedMsg, err := l.convertCtypeMsgToGrpcMsg(msg)
	if err != nil {
		fmt.Println("Error converting msg:", err)
		return nil, err
	}
	err = l.svcCtx.DB.Preload("SendUserModel").First(&chatModel, chatModel.Id).Error
	if err != nil {
		fmt.Println("preload异常", err.Error())
		return nil, err
	}

	return &chat_rpc.SendMsgRes{
		Id:             uint32(chatModel.Id),
		MessageId:      chatModel.MessageID, // 支持 uint32 类型
		ConversationId: chatModel.ConversationID,
		Msg:            convertedMsg,
		MsgPreview:     chatModel.MsgPreview,
		Sender: &chat_rpc.Sender{
			UserId:   chatModel.SendUserModel.UUID,
			Avatar:   chatModel.SendUserModel.Avatar,
			Nickname: chatModel.SendUserModel.NickName,
		},
		CreateAt: chatModel.CreatedAt.String(),
		Status:   1,             // 1:正常状态
		Seq:      chatModel.Seq, // 消息序列号
	}, nil
}

func (l *SendMsgLogic) convertCtypeMsgToGrpcMsg(msg ctype.Msg) (*chat_rpc.Msg, error) {
	// 将 ctype.Msg 转换为 JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 将 JSON 解析为 chat_rpc.Msg
	var convertedMsg chat_rpc.Msg
	err = json.Unmarshal(jsonData, &convertedMsg)
	if err != nil {
		return nil, err
	}

	return &convertedMsg, nil
}
