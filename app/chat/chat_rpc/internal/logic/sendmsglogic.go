package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/chat/chat_utils"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/ajax"
	"beaver/common/models/ctype"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
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

	// 获取下一个消息序列号（消息表内部序列号）
	var maxSeq int64
	err := l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("conversation_id = ?", in.ConversationId).
		Select("COALESCE(MAX(seq), 0)").
		Scan(&maxSeq).Error
	if err != nil {
		l.Logger.Errorf("获取消息序列号失败: conversationId=%s, error=%v", in.ConversationId, err)
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

	// 1. 创建消息记录并设置预览
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	err = l.svcCtx.DB.Create(&chatModel).Error
	if err != nil {
		l.Logger.Errorf("创建消息记录失败: conversationId=%s, userId=%s, error=%v", in.ConversationId, in.UserId, err)
		return nil, err
	}

	// 1.1 构建 Messages 表更新数据
	messagesUpdate := map[string]interface{}{
		"table":          "messages",
		"conversationId": in.ConversationId,
		"data": []map[string]interface{}{
			{
				"seq": chatModel.Seq,
			},
		},
	}

	// 2. 更新会话级别的信息
	conversationVersion, err := chat_utils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversationType, chatModel.Seq, chatModel.MsgPreview)
	if err != nil {
		l.Logger.Errorf("更新会话信息失败: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, err
	}

	// 2.1 构建 Conversations 表更新数据
	conversationsUpdate := map[string]interface{}{
		"table":          "conversations",
		"conversationId": in.ConversationId,
		"data": []map[string]interface{}{
			{
				"version": int32(conversationVersion),
			},
		},
	}

	// 3. 批量更新该会话所有用户的会话关系（包括发送者：恢复隐藏状态，更新版本号，更新已读序列号）
	allUserConversationUpdates, err := chat_utils.UpdateAllUserConversationsInChat(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, in.UserId, chatModel.Seq)
	if err != nil {
		l.Logger.Errorf("批量更新用户会话关系失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 不影响消息发送成功，只记录错误
		allUserConversationUpdates = []chat_utils.UserConversationUpdate{}
	}

	// 转换消息格式
	convertedMsg, err := l.convertCtypeMsgToGrpcMsg(msg)
	if err != nil {
		l.Logger.Errorf("转换消息格式失败: %v", err)
		return nil, err
	}

	// 获取发送者信息
	sender, err := l.getSenderInfo(chatModel)
	if err != nil {
		l.Logger.Errorf("获取发送者信息失败: %v", err)
		return nil, err
	}

	// 5. 异步推送消息更新给会话成员（按表分组，一次推送包含所有相关更新）
	go func() {
		l.notifyMessageUpdateGrouped(in.ConversationId, in.UserId, conversationType, messagesUpdate, conversationsUpdate, allUserConversationUpdates)
	}()

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

// notifyMessageUpdateGrouped 按会话分组推送消息更新（给该会话的所有用户推送）
func (l *SendMsgLogic) notifyMessageUpdateGrouped(conversationId, senderId string, conversationType int, messagesUpdate, conversationsUpdate map[string]interface{}, allUserConversationUpdates []chat_utils.UserConversationUpdate) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("推送消息更新时发生panic: %v", r)
		}
	}()

	// 获取该会话的所有用户ID（包括发送者，给所有人推送）
	var recipientIds []string
	for _, update := range allUserConversationUpdates {
		if update.ConversationID == conversationId {
			recipientIds = append(recipientIds, update.UserID)
		}
	}

	// 构建按表分组的更新数据结构
	tableUpdates := []map[string]interface{}{
		messagesUpdate,      // 已经构建好的 Messages 更新
		conversationsUpdate, // 已经构建好的 Conversations 更新
	}

	// 3. User_conversations 表更新（直接循环 allUserConversationUpdates）
	for _, update := range allUserConversationUpdates {
		userConversationsUpdate := map[string]interface{}{
			"table":          "user_conversations",
			"userId":         update.UserID,
			"conversationId": update.ConversationID,
			"data": []map[string]interface{}{
				{
					"version": int32(update.Version),
				},
			},
		}
		tableUpdates = append(tableUpdates, userConversationsUpdate)
	}

	// 为每个接收者推送批量更新
	messageType := wsTypeConst.ChatConversationMessageReceive

	fmt.Println(("111111111111111111111"))
	fmt.Println(("111111111111111111111"))
	fmt.Println(("111111111111111111111"))
	fmt.Println(("111111111111111111111"))
	fmt.Println(("111111111111111111111"))
	fmt.Println("推送消息更新给会话成员: ", recipientIds)
	for _, recipientId := range recipientIds {
		// 一次性推送所有表的更新信息
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd.Hosts[0], wsCommandConst.CHAT_MESSAGE, messageType, senderId, recipientId, map[string]interface{}{
			"tableUpdates": tableUpdates, // 按表分组的更新数组
		}, conversationId)

		l.Logger.Infof("分组推送消息相关更新: recipient=%s, conversation=%s, tableUpdateCount=%d",
			recipientId, conversationId, len(tableUpdates))
	}
}

// getSenderInfo 获取发送者信息
func (l *SendMsgLogic) getSenderInfo(chatModel chat_models.ChatMessage) (*chat_rpc.Sender, error) {
	sendUserID := ""
	if chatModel.SendUserID != nil {
		sendUserID = *chatModel.SendUserID
	}

	if sendUserID == "" {
		// 通知消息：SendUserID为空
		return &chat_rpc.Sender{
			UserId:   "",
			NickName: "通知消息",
			Avatar:   "",
		}, nil
	}

	// 调用UserRpc获取用户信息
	userInfoResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: sendUserID,
	})
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: userId=%s, error=%v", sendUserID, err)
		// 返回默认发送者信息，不影响消息发送
		return &chat_rpc.Sender{
			UserId:   sendUserID,
			NickName: "未知用户",
			Avatar:   "",
		}, nil
	}

	// 使用结构化用户信息
	userInfo := userInfoResp.UserInfo
	return &chat_rpc.Sender{
		UserId:   sendUserID,
		NickName: userInfo.NickName,
		Avatar:   userInfo.Avatar,
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
