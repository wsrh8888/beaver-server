package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
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
		if !strings.Contains(in.ConversationId, in.UserID) {
			logx.Errorf("用户id不匹配，用户id：%s，会话id：%s", in.UserID, in.ConversationId)
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
				FileId: in.Msg.ImageMsg.FileId,
				Name:   in.Msg.ImageMsg.Name,
			},
		}
	case ctype.VideoMsgType:
		msg = ctype.Msg{
			Type: ctype.VideoMsgType,
			VideoMsg: &ctype.VideoMsg{
				Src:   in.Msg.VideoMsg.Src,
				Title: in.Msg.VideoMsg.Title,
				Time:  in.Msg.VideoMsg.Time,
			},
		}
	case ctype.FileMsgType:
		msg = ctype.Msg{
			Type: ctype.FileMsgType,
			FileMsg: &ctype.FileMsg{
				Src:   in.Msg.FileMsg.Src,
				Title: in.Msg.FileMsg.Title,
				Size:  in.Msg.FileMsg.Size,
				Type:  in.Msg.FileMsg.Type,
			},
		}
	case ctype.VoiceMsgType:
		msg = ctype.Msg{
			Type: ctype.VoiceMsgType,
			VoiceMsg: &ctype.VoiceMsg{
				Src:  in.Msg.VoiceMsg.Src,
				Time: in.Msg.VoiceMsg.Time,
			},
		}
	default:
		return nil, fmt.Errorf("未识别到该类型: %d", msgType)
	}

	chatModel := chat_models.ChatModel{
		SendUserID:     in.UserID,
		ConversationID: in.ConversationId,
		MsgType:        msgType,
		Msg:            &msg,
	}
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()

	err := l.svcCtx.DB.Create(&chatModel).Preload("SendUserModel").Error
	if err != nil {
		return nil, err
	}

	err = l.updateUserConversations(in.ConversationId, in.UserID, chatModel.MsgPreview)
	if err != nil {
		return nil, err
	}

	convertedMsg, err := l.convertCtypeMsgToGrpcMsg(msg)
	if err != nil {
		fmt.Println("Error converting msg:", err)
		return nil, err
	}
	err = l.svcCtx.DB.Preload("SendUserModel").First(&chatModel, chatModel.ID).Error
	if err != nil {
		fmt.Println("preload异常", err.Error())
		return nil, err
	}

	return &chat_rpc.SendMsgRes{
		MessageId:      uint32(chatModel.ID), // 支持 uint32 类型
		ConversationId: chatModel.ConversationID,
		Msg:            convertedMsg,
		MsgPreview:     chatModel.MsgPreview,
		Sender: &chat_rpc.Sender{
			UserID:   chatModel.SendUserModel.UUID,
			Avatar:   chatModel.SendUserModel.Avatar,
			Nickname: chatModel.SendUserModel.NickName,
		},
		CreateAt: chatModel.CreatedAt.String(),
	}, nil
}

func (l *SendMsgLogic) updateUserConversations(conversationID, userID, lastMessage string) error {
	var userConvo chat_models.ChatUserConversationModel
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&userConvo).Error
	if err != nil {
		if err := l.svcCtx.DB.Create(&chat_models.ChatUserConversationModel{
			UserID:         userID,
			ConversationID: conversationID,
			LastMessage:    lastMessage,
			IsDeleted:      false,
		}).Error; err != nil {
			return err
		}
	} else {
		if err := l.svcCtx.DB.Model(&userConvo).
			Updates(map[string]interface{}{
				"last_message": lastMessage,
				"is_deleted":   false,
				"updated_at":   time.Now(), // 添加这一行
			}).Error; err != nil {
			return err
		}
	}
	return nil
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
