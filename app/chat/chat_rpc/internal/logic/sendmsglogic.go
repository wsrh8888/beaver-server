package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	chatrpcutils "beaver/app/chat/chat_rpc/internal/utils"
	"beaver/app/friend/friend_models"
	"beaver/app/group/group_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
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

// BuildMsgFromProto 将 Proto 的消息结构转换/映射为 ctype.Msg 逻辑（支持回复模式递归）
func (l *SendMsgLogic) BuildMsgFromProto(protoMsg *chat_rpc.Msg) *ctype.Msg {
	if protoMsg == nil {
		return nil
	}

	msgType := ctype.MsgType(protoMsg.Type)
	var msg ctype.Msg

	switch msgType {
	case ctype.TextMsgType:
		if protoMsg.TextMsg != nil {
			msg = ctype.Msg{
				Type: ctype.TextMsgType,
				TextMsg: &ctype.TextMsg{
					Content: protoMsg.TextMsg.Content,
				},
			}
		}
	case ctype.ImageMsgType:
		if protoMsg.ImageMsg != nil {
			imageMsg := &ctype.ImageMsg{
				FileKey: protoMsg.ImageMsg.FileKey,
			}
			if protoMsg.ImageMsg.Width > 0 {
				imageMsg.Width = int(protoMsg.ImageMsg.Width)
			}
			if protoMsg.ImageMsg.Height > 0 {
				imageMsg.Height = int(protoMsg.ImageMsg.Height)
			}
			if protoMsg.ImageMsg.Size > 0 {
				imageMsg.Size = protoMsg.ImageMsg.Size
			}
			msg = ctype.Msg{
				Type:     ctype.ImageMsgType,
				ImageMsg: imageMsg,
			}
		}
	case ctype.VideoMsgType:
		if protoMsg.VideoMsg != nil {
			videoMsg := &ctype.VideoMsg{
				FileKey: protoMsg.VideoMsg.FileKey,
			}
			if protoMsg.VideoMsg.Width > 0 {
				videoMsg.Width = int(protoMsg.VideoMsg.Width)
			}
			if protoMsg.VideoMsg.Height > 0 {
				videoMsg.Height = int(protoMsg.VideoMsg.Height)
			}
			if protoMsg.VideoMsg.Duration > 0 {
				videoMsg.Duration = int(protoMsg.VideoMsg.Duration)
			}
			if protoMsg.VideoMsg.ThumbnailKey != "" {
				videoMsg.ThumbnailKey = protoMsg.VideoMsg.ThumbnailKey
			}
			if protoMsg.VideoMsg.Size > 0 {
				videoMsg.Size = protoMsg.VideoMsg.Size
			}
			msg = ctype.Msg{
				Type:     ctype.VideoMsgType,
				VideoMsg: videoMsg,
			}
		}
	case ctype.FileMsgType:
		if protoMsg.FileMsg != nil {
			fileMsg := &ctype.FileMsg{
				FileKey: protoMsg.FileMsg.FileKey,
			}
			if protoMsg.FileMsg.FileName != "" {
				fileMsg.FileName = protoMsg.FileMsg.FileName
			}
			if protoMsg.FileMsg.Size > 0 {
				fileMsg.Size = protoMsg.FileMsg.Size
			}
			if protoMsg.FileMsg.MimeType != "" {
				fileMsg.MimeType = protoMsg.FileMsg.MimeType
			}
			if protoMsg.FileMsg.Extension != "" {
				fileMsg.Extension = protoMsg.FileMsg.Extension
			}
			if protoMsg.FileMsg.OpenMode > 0 {
				fileMsg.OpenMode = int(protoMsg.FileMsg.OpenMode)
			}
			msg = ctype.Msg{
				Type:    ctype.FileMsgType,
				FileMsg: fileMsg,
			}
		}
	case ctype.VoiceMsgType:
		if protoMsg.VoiceMsg != nil {
			voiceMsg := &ctype.VoiceMsg{
				FileKey: protoMsg.VoiceMsg.FileKey,
			}
			if protoMsg.VoiceMsg.Duration > 0 {
				voiceMsg.Duration = int(protoMsg.VoiceMsg.Duration)
			}
			if protoMsg.VoiceMsg.Size > 0 {
				voiceMsg.Size = protoMsg.VoiceMsg.Size
			}
			msg = ctype.Msg{
				Type:     ctype.VoiceMsgType,
				VoiceMsg: voiceMsg,
			}
		}
	case ctype.EmojiMsgType:
		if protoMsg.EmojiMsg != nil {
			msg = ctype.Msg{
				Type: ctype.EmojiMsgType,
				EmojiMsg: &ctype.EmojiMsg{
					FileKey:   protoMsg.EmojiMsg.FileKey,
					EmojiID:   protoMsg.EmojiMsg.EmojiId,
					PackageID: protoMsg.EmojiMsg.PackageId,
					Width:     protoMsg.EmojiMsg.Width,
					Height:    protoMsg.EmojiMsg.Height,
				},
			}
		}
	case ctype.NotificationMsgType:
		if protoMsg.NotificationMsg != nil {
			notificationMsg := &ctype.NotificationMsg{
				Type:   int(protoMsg.NotificationMsg.Type),
				Actors: protoMsg.NotificationMsg.Actors,
			}
			msg = ctype.Msg{
				Type:            ctype.NotificationMsgType,
				NotificationMsg: notificationMsg,
			}
		}
	case ctype.AudioFileMsgType:
		if protoMsg.AudioFileMsg != nil {
			audioFileMsg := &ctype.AudioFileMsg{
				FileKey: protoMsg.AudioFileMsg.FileKey,
			}
			if protoMsg.AudioFileMsg.FileName != "" {
				audioFileMsg.FileName = protoMsg.AudioFileMsg.FileName
			}
			if protoMsg.AudioFileMsg.Duration > 0 {
				audioFileMsg.Duration = int(protoMsg.AudioFileMsg.Duration)
			}
			if protoMsg.AudioFileMsg.Size > 0 {
				audioFileMsg.Size = protoMsg.AudioFileMsg.Size
			}
			msg = ctype.Msg{
				Type:         ctype.AudioFileMsgType,
				AudioFileMsg: audioFileMsg,
			}
		}
	case ctype.CallMsgType:
		if protoMsg.CallMsg != nil {
			msg = ctype.Msg{
				Type: ctype.CallMsgType,
				CallMsg: &ctype.CallMsg{
					RoomID:   protoMsg.CallMsg.RoomId,
					CallType: int(protoMsg.CallMsg.CallType),
					Status:   int(protoMsg.CallMsg.Status),
					Duration: protoMsg.CallMsg.Duration,
				},
			}
		}
	case ctype.WithdrawMsgType:
		if protoMsg.WithdrawMsg != nil {
			msg = ctype.Msg{
				Type: ctype.WithdrawMsgType,
				WithdrawMsg: &ctype.WithdrawMsg{
					OriginMsgID: protoMsg.WithdrawMsg.OriginMsgId,
					OriginMsg:   l.BuildMsgFromProto(protoMsg.WithdrawMsg.OriginMsg), // 支持快照递归
				},
			}
		}
	case ctype.ReplyMsgType:
		if protoMsg.ReplyMsg != nil {
			// 这里递归转换被回复的消息 (引用部分)
			var originMsg *ctype.Msg
			if protoMsg.ReplyMsg.OriginMsg != nil {
				originMsg = l.BuildMsgFromProto(protoMsg.ReplyMsg.OriginMsg)
			}

			// 这里递归转换回复的具体内容 (回复内容部分)
			var replyMsg *ctype.Msg
			if protoMsg.ReplyMsg.ReplyMsg != nil {
				replyMsg = l.BuildMsgFromProto(protoMsg.ReplyMsg.ReplyMsg)
			}

			msg = ctype.Msg{
				Type: ctype.ReplyMsgType,
				ReplyMsg: &ctype.ReplyMsg{
					OriginMsgID: protoMsg.ReplyMsg.OriginMsgId,
					OriginMsg:   originMsg,
					ReplyMsg:    replyMsg,
				},
			}
		}
	case ctype.ForwardMsgType:
		if protoMsg.ForwardMsg != nil {
			msg = ctype.Msg{
				Type: ctype.ForwardMsgType,
				ForwardMsg: &ctype.ForwardMsg{
					Title:    protoMsg.ForwardMsg.Title,
					RecordID: protoMsg.ForwardMsg.RecordId,
					Count:    int(protoMsg.ForwardMsg.Count),
				},
			}
		}
	case ctype.MarkdownMsgType:
		if protoMsg.MarkdownMsg != nil {
			msg = ctype.Msg{
				Type: ctype.MarkdownMsgType,
				MarkdownMsg: &ctype.MarkdownMsg{
					Content: protoMsg.MarkdownMsg.Content,
					Title:   protoMsg.MarkdownMsg.Title,
				},
			}
		}
	case ctype.LinkMsgType:
		if protoMsg.LinkMsg != nil {
			msg = ctype.Msg{
				Type: ctype.LinkMsgType,
				LinkMsg: &ctype.LinkMsg{
					URL:      protoMsg.LinkMsg.Url,
					Title:    protoMsg.LinkMsg.Title,
					Desc:     protoMsg.LinkMsg.Desc,
					ImageURL: protoMsg.LinkMsg.ImageUrl,
				},
			}
		}
	case ctype.CloudDocMsgType:
		if protoMsg.CloudDocMsg != nil {
			msg = ctype.Msg{
				Type: ctype.CloudDocMsgType,
				CloudDocMsg: &ctype.CloudDocMsg{
					DocID:    protoMsg.CloudDocMsg.DocId,
					DocType:  int(protoMsg.CloudDocMsg.DocType),
					Title:    protoMsg.CloudDocMsg.Title,
					OwnerID:  protoMsg.CloudDocMsg.OwnerId,
					Perm:     int(protoMsg.CloudDocMsg.Perm),
					CoverURL: protoMsg.CloudDocMsg.CoverUrl,
					Revision: protoMsg.CloudDocMsg.Revision,
				},
			}
		}
	}
	if len(protoMsg.AtUserIds) > 0 {
		msg.AtUserIDs = protoMsg.AtUserIds
	}
	return &msg
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

		var friendModel friend_models.FriendModel
		if !friendModel.IsFriend(l.svcCtx.DB, userIds[0], userIds[1]) {
			logx.Errorf("不是好友关系，用户IDs: %v", userIds)
			return nil, errors.New("不是好友关系")
		}

		// 黑名单检查（双向）：任一方拉黑对方均无法发消息
		var blockCount int64
		l.svcCtx.DB.Model(&friend_models.FriendBlockModel{}).
			Where("(user_id = ? AND blocked_user_id = ?) OR (user_id = ? AND blocked_user_id = ?)",
				userIds[0], userIds[1], userIds[1], userIds[0]).
			Count(&blockCount)
		if blockCount > 0 {
			return nil, errors.New("无法发送消息")
		}
	} else if conversationType == 2 {
		// 群聊：禁言检查
		groupID := conversation.GetTargetIDByConversation(in.ConversationId, in.UserId)

		var member group_models.GroupMemberModel
		if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", groupID, in.UserId).
			First(&member).Error; err != nil {
			return nil, errors.New("你不是该群成员")
		}

		// 群主(1)和管理员(2)不受禁言限制；群通知机器人（nbot_前缀）也跳过禁言校验
		if member.Role != 1 && member.Role != 2 && !strings.HasPrefix(in.UserId, "nbot_") {
			// 全员禁言检查
			var group group_models.GroupModel
			if err := l.svcCtx.DB.Where("group_id = ?", groupID).First(&group).Error; err == nil {
				if group.IsMuteAll {
					return nil, errors.New("当前群已开启全员禁言")
				}
			}
			// 个人禁言检查
			if member.MutedUntil != nil && member.MutedUntil.After(time.Now()) {
				return nil, errors.New("你已被禁言")
			}
		}
	}

	// 调用抽离好的递归转换函数
	msg := l.BuildMsgFromProto(in.Msg)
	msgType := msg.Type

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
		Msg:              msg,
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
	conversationVersion, err := chatrpcutils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversationType, chatModel.Seq, chatModel.MsgPreview)
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
	allUserConversationUpdates, err := chatrpcutils.UpdateAllUserConversationsInChat(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, in.UserId, chatModel.Seq)
	if err != nil {
		l.Logger.Errorf("批量更新用户会话关系失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 不影响消息发送成功，只记录错误
		allUserConversationUpdates = []chatrpcutils.UserConversationUpdate{}
	}

	// 转换消息格式
	convertedMsg, err := l.convertCtypeMsgToGrpcMsg(*msg)
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
		l.sendOfflinePushIfNeeded(in.ConversationId, in.UserId, chatModel, recipientIdsFromUpdates(allUserConversationUpdates, in.ConversationId))
	}()

	return &chat_rpc.SendMsgRes{
		Id:               uint32(chatModel.Id),
		MessageId:        chatModel.MessageID,
		ConversationId:   chatModel.ConversationID,
		Msg:              convertedMsg,
		MsgPreview:       chatModel.MsgPreview,
		Sender:           sender,
		CreatedAt:        chatModel.CreatedAt.String(),
		Status:           1,
		ConversationType: uint32(chatModel.ConversationType),
		Seq:              chatModel.Seq,
	}, nil
}

// notifyMessageUpdateGrouped 按会话分组推送消息更新（给该会话的所有用户推送）
func (l *SendMsgLogic) notifyMessageUpdateGrouped(conversationId, senderId string, conversationType int, messagesUpdate, conversationsUpdate map[string]interface{}, allUserConversationUpdates []chatrpcutils.UserConversationUpdate) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("推送消息更新时发生panic: %v", r)
		}
	}()

	var recipientIds []string
	for _, update := range allUserConversationUpdates {
		if update.ConversationID == conversationId {
			recipientIds = append(recipientIds, update.UserID)
		}
	}

	// 构建 userId -> user_conversations 更新的快速查找 map
	userConvUpdateMap := make(map[string]map[string]interface{}, len(allUserConversationUpdates))
	for _, update := range allUserConversationUpdates {
		userConvUpdateMap[update.UserID] = map[string]interface{}{
			"table":          "user_conversations",
			"userId":         update.UserID,
			"conversationId": update.ConversationID,
			"data": []map[string]interface{}{
				{
					"version": int32(update.Version),
				},
			},
		}
	}

	// 为每个接收者推送：只携带自己的 user_conversations，不携带其他成员的
	messageType := wsTypeConst.ChatConversationMessageReceive

	for _, recipientId := range recipientIds {
		tableUpdates := []map[string]interface{}{
			messagesUpdate,
			conversationsUpdate,
		}
		if uc, ok := userConvUpdateMap[recipientId]; ok {
			tableUpdates = append(tableUpdates, uc)
		}

		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     messageType,
			"senderId": senderId,
			"targetId": recipientId,
			"body": map[string]interface{}{
				"tableUpdates": tableUpdates,
			},
			"conversationId": conversationId,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
			l.Logger.Errorf("MQ 推送失败: recipient=%s, conversation=%s, error=%v", recipientId, conversationId, err)
		}
	}
}

func recipientIdsFromUpdates(updates []chatrpcutils.UserConversationUpdate, conversationID string) []string {
	var ids []string
	for _, update := range updates {
		if update.ConversationID == conversationID {
			ids = append(ids, update.UserID)
		}
	}
	return ids
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
		return &chat_rpc.Sender{
			UserId:   sendUserID,
			NickName: "未知用户",
			Avatar:   "",
		}, nil
	}

	userInfo := userInfoResp.UserInfo
	return &chat_rpc.Sender{
		UserId:   sendUserID,
		NickName: userInfo.NickName,
		Avatar:   userInfo.Avatar,
	}, nil
}

func (l *SendMsgLogic) convertCtypeMsgToGrpcMsg(m ctype.Msg) (*chat_rpc.Msg, error) {
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
			convertedOrigin, _ := l.convertCtypeMsgToGrpcMsg(*m.WithdrawMsg.OriginMsg)
			rpcMsg.WithdrawMsg = &chat_rpc.WithdrawMsg{
				OriginMsgId: m.WithdrawMsg.OriginMsgID,
				OriginMsg:   convertedOrigin,
			}
		}
	case ctype.ReplyMsgType:
		if m.ReplyMsg != nil {
			var convertedOrigin *chat_rpc.Msg
			if m.ReplyMsg.OriginMsg != nil {
				convertedOrigin, _ = l.convertCtypeMsgToGrpcMsg(*m.ReplyMsg.OriginMsg)
			}
			var convertedReply *chat_rpc.Msg
			if m.ReplyMsg.ReplyMsg != nil {
				convertedReply, _ = l.convertCtypeMsgToGrpcMsg(*m.ReplyMsg.ReplyMsg)
			}
			rpcMsg.ReplyMsg = &chat_rpc.ReplyMsg{
				OriginMsgId: m.ReplyMsg.OriginMsgID,
				OriginMsg:   convertedOrigin,
				ReplyMsg:    convertedReply,
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
				DocId:    m.CloudDocMsg.DocID,
				DocType:  int32(m.CloudDocMsg.DocType),
				Title:    m.CloudDocMsg.Title,
				OwnerId:  m.CloudDocMsg.OwnerID,
				Perm:     int32(m.CloudDocMsg.Perm),
				CoverUrl: m.CloudDocMsg.CoverURL,
				Revision: m.CloudDocMsg.Revision,
			}
		}
	}
	return rpcMsg, nil
}
