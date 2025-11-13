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
	"beaver/app/user/user_rpc/types/user_rpc"
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
	conversationType, userIds := conversation.ParseConversationWithType(in.ConversationId)

	if conversationType == 1 {
		// 私聊需要验证好友关系
		if !strings.Contains(in.ConversationId, in.UserId) {
			logx.Errorf("用户id不匹配，用户id：%s，会话id：%s", in.UserId, in.ConversationId)
			return nil, errors.New("异常操作")
		}

		if len(userIds) != 2 {
			logx.Errorf("私聊会话ID解析失败，期望2个用户ID，实际: %v", userIds)
			return nil, errors.New("无效的私聊会话ID")
		}

		var friend friend_models.FriendModel
		if !friend.IsFriend(l.svcCtx.DB, userIds[0], userIds[1]) {
			logx.Errorf("不是好友关系，用户IDs: %v", userIds)
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
		imageMsg := &ctype.ImageMsg{
			FileKey: in.Msg.ImageMsg.FileKey,
		}
		// 设置可选字段（proto 字段名是小写开头，但 Go 生成的是大写开头）
		// 注意：需要重新生成 proto 后，字段名会是 Width, Height, Size
		if in.Msg.ImageMsg.Width > 0 {
			imageMsg.Width = int(in.Msg.ImageMsg.Width)
		}
		if in.Msg.ImageMsg.Height > 0 {
			imageMsg.Height = int(in.Msg.ImageMsg.Height)
		}
		if in.Msg.ImageMsg.Size > 0 {
			imageMsg.Size = in.Msg.ImageMsg.Size
		}
		msg = ctype.Msg{
			Type:     ctype.ImageMsgType,
			ImageMsg: imageMsg,
		}
	case ctype.VideoMsgType:
		videoMsg := &ctype.VideoMsg{
			FileKey: in.Msg.VideoMsg.FileKey,
		}
		// 设置可选字段
		if in.Msg.VideoMsg.Width > 0 {
			videoMsg.Width = int(in.Msg.VideoMsg.Width)
		}
		if in.Msg.VideoMsg.Height > 0 {
			videoMsg.Height = int(in.Msg.VideoMsg.Height)
		}
		if in.Msg.VideoMsg.Duration > 0 {
			videoMsg.Duration = int(in.Msg.VideoMsg.Duration)
		}
		if in.Msg.VideoMsg.ThumbnailKey != "" {
			videoMsg.ThumbnailKey = in.Msg.VideoMsg.ThumbnailKey
		}
		if in.Msg.VideoMsg.Size > 0 {
			videoMsg.Size = in.Msg.VideoMsg.Size
		}
		msg = ctype.Msg{
			Type:     ctype.VideoMsgType,
			VideoMsg: videoMsg,
		}
	case ctype.FileMsgType:
		fileMsg := &ctype.FileMsg{
			FileKey: in.Msg.FileMsg.FileKey,
		}
		// 设置可选字段
		if in.Msg.FileMsg.FileName != "" {
			fileMsg.FileName = in.Msg.FileMsg.FileName
		}
		if in.Msg.FileMsg.Size > 0 {
			fileMsg.Size = in.Msg.FileMsg.Size
		}
		if in.Msg.FileMsg.MimeType != "" {
			fileMsg.MimeType = in.Msg.FileMsg.MimeType
		}
		msg = ctype.Msg{
			Type:    ctype.FileMsgType,
			FileMsg: fileMsg,
		}
	case ctype.VoiceMsgType:
		voiceMsg := &ctype.VoiceMsg{
			FileKey: in.Msg.VoiceMsg.FileKey,
		}
		// 设置可选字段
		if in.Msg.VoiceMsg.Duration > 0 {
			voiceMsg.Duration = int(in.Msg.VoiceMsg.Duration)
		}
		if in.Msg.VoiceMsg.Size > 0 {
			voiceMsg.Size = in.Msg.VoiceMsg.Size
		}
		msg = ctype.Msg{
			Type:     ctype.VoiceMsgType,
			VoiceMsg: voiceMsg,
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
	case ctype.NotificationMsgType:
		notificationMsg := &ctype.NotificationMsg{
			Type:   int(in.Msg.NotificationMsg.Type),
			Actors: in.Msg.NotificationMsg.Actors,
		}
		msg = ctype.Msg{
			Type:            ctype.NotificationMsgType,
			NotificationMsg: notificationMsg,
		}
	case ctype.AudioFileMsgType:
		audioFileMsg := &ctype.AudioFileMsg{
			FileKey: in.Msg.AudioFileMsg.FileKey,
		}
		// 设置可选字段
		if in.Msg.AudioFileMsg.FileName != "" {
			audioFileMsg.FileName = in.Msg.AudioFileMsg.FileName
		}
		if in.Msg.AudioFileMsg.Duration > 0 {
			audioFileMsg.Duration = int(in.Msg.AudioFileMsg.Duration)
		}
		if in.Msg.AudioFileMsg.Size > 0 {
			audioFileMsg.Size = in.Msg.AudioFileMsg.Size
		}
		msg = ctype.Msg{
			Type:         ctype.AudioFileMsgType,
			AudioFileMsg: audioFileMsg,
		}
	default:
		return nil, fmt.Errorf("未识别到该类型: %d", msgType)
	}

	// 获取下一个序列号
	var maxSeq int64
	err := l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("conversation_id = ?", in.ConversationId).
		Select("COALESCE(MAX(seq), 0)").
		Scan(&maxSeq).Error
	if err != nil {
		l.Logger.Errorf("获取序列号失败: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, err
	}

	nextSeq := maxSeq + 1

	chatModel := chat_models.ChatMessage{
		SendUserID:       &in.UserId,
		MessageID:        in.MessageId,
		ConversationID:   in.ConversationId,
		ConversationType: conversationType,
		Seq:              nextSeq, // 设置正确的序列号
		MsgType:          msgType,
		Msg:              &msg,
	}

	// 1. 创建消息记录
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	err = l.svcCtx.DB.Create(&chatModel).Error
	if err != nil {
		return nil, err
	}

	// 2. 更新会话级别的信息
	err = chat_utils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversationType, chatModel.Seq, chatModel.MsgPreview)
	if err != nil {
		return nil, err
	}

	// 3. 更新用户会话关系
	err = chat_utils.UpdateUserConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.UserId, in.ConversationId, false)
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

	// 重新查询消息（不使用Preload）
	err = l.svcCtx.DB.First(&chatModel, chatModel.Id).Error
	if err != nil {
		fmt.Println("查询消息异常", err.Error())
		return nil, err
	}

	// 获取发送者信息
	var sender *chat_rpc.Sender
	sendUserID := ""
	if chatModel.SendUserID != nil {
		sendUserID = *chatModel.SendUserID
	}

	if sendUserID != "" {
		// 调用UserRpc获取用户信息
		userInfoResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
			UserID: sendUserID,
		})
		if err != nil {
			l.Logger.Errorf("获取用户信息失败: %v", err)
			// 设置默认发送者信息
			sender = &chat_rpc.Sender{
				UserId:   sendUserID,
				Nickname: "未知用户",
				Avatar:   "",
			}
		} else {
			// 直接使用结构化用户信息
			userInfo := userInfoResp.UserInfo
			sender = &chat_rpc.Sender{
				UserId:   sendUserID,
				Nickname: userInfo.NickName,
				Avatar:   userInfo.Avatar,
			}
		}
	} else {
		// 通知消息：SendUserID为空
		sender = &chat_rpc.Sender{
			UserId:   "",
			Nickname: "通知消息",
			Avatar:   "",
		}
	}

	return &chat_rpc.SendMsgRes{
		Id:               uint32(chatModel.Id),
		MessageId:        chatModel.MessageID,
		ConversationId:   chatModel.ConversationID,
		Msg:              convertedMsg,
		MsgPreview:       chatModel.MsgPreview,
		Sender:           sender,
		CreateAt:         chatModel.CreatedAt.String(),
		Status:           1,
		ConversationType: uint32(chatModel.ConversationType),
		Seq:              chatModel.Seq,
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
