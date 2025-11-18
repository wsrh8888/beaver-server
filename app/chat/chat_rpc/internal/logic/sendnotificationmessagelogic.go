package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/chat/chat_utils"
	"beaver/common/ajax"
	"beaver/common/models/ctype"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendNotificationMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendNotificationMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNotificationMessageLogic {
	return &SendNotificationMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendNotificationMessageLogic) SendNotificationMessage(in *chat_rpc.SendNotificationMessageReq) (*chat_rpc.SendNotificationMessageRes, error) {
	// 参数验证
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}

	if in.Content == "" {
		return nil, errors.New("消息内容不能为空")
	}

	// 构建通知消息的结构化数据
	actors := []string{}
	if in.RelatedUserId != "" {
		actors = append(actors, in.RelatedUserId)
	}

	// 检查会话是否存在
	var conversation chat_models.ChatConversationMeta
	err := l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).First(&conversation).Error
	if err != nil {
		l.Logger.Errorf("会话不存在: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("会话不存在")
	}

	// 生成消息ID
	messageId := uuid.New().String()

	// 获取下一个序列号
	var maxSeq int64
	err = l.svcCtx.DB.Model(&chat_models.ChatMessage{}).
		Where("conversation_id = ?", in.ConversationId).
		Select("COALESCE(MAX(seq), 0)").
		Scan(&maxSeq).Error
	if err != nil {
		l.Logger.Errorf("获取序列号失败: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("获取序列号失败")
	}

	nextSeq := maxSeq + 1

	// 为通知消息创建Msg结构（使用NotificationMsg类型）
	notificationMsg := &ctype.Msg{
		Type: ctype.NotificationMsgType, // 使用通知消息类型
		NotificationMsg: &ctype.NotificationMsg{
			Type:   int(in.MessageType), // 通知类型
			Actors: actors,              // 相关用户ID列表
		},
	}

	notificationMessage := chat_models.ChatMessage{
		MessageID:        messageId,
		ConversationID:   in.ConversationId,
		ConversationType: conversation.Type,
		Seq:              nextSeq,
		SendUserID:       nil, // 通知消息SendUserID为null
		MsgType:          7,   // 通知消息类型
		MsgPreview:       in.Content,
		Msg:              notificationMsg, // 通知消息的结构化内容
	}

	if err := l.svcCtx.DB.Create(&notificationMessage).Error; err != nil {
		l.Logger.Errorf("创建通知消息失败: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("创建通知消息失败")
	}

	// 更新会话级别的信息，获取会话版本号
	conversationVersion, err := chat_utils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversation.Type, nextSeq, in.Content)
	if err != nil {
		l.Logger.Errorf("更新会话信息失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 这里不返回错误，因为消息已经创建成功
		conversationVersion = 0
	}

	// 构建会话更新数据
	conversationsUpdate := map[string]interface{}{
		"table":          "conversations",
		"conversationId": in.ConversationId,
		"data": []map[string]interface{}{
			{
				"version": int32(conversationVersion),
			},
		},
	}

	// 批量更新该会话所有用户的会话关系（恢复隐藏状态，更新版本号）
	// 注意：通知消息没有发送者，所以不处理已读序列号
	allUserConversationUpdates, err := chat_utils.UpdateAllUserConversationsInChat(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, "", 0)
	if err != nil {
		l.Logger.Errorf("批量更新用户会话关系失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 不影响消息发送成功，只记录错误
		allUserConversationUpdates = []chat_utils.UserConversationUpdate{}
	}

	// 如果有发送者（RelatedUserId），则推送通知消息给所有相关用户
	if in.RelatedUserId != "" {
		// 构建 Messages 表更新数据
		messagesUpdate := map[string]interface{}{
			"table":          "messages",
			"conversationId": in.ConversationId,
			"data": []map[string]interface{}{
				{
					"seq": nextSeq,
				},
			},
		}

		// 异步推送消息更新给会话成员
		go func() {
			l.notifyNotificationUpdateGrouped(in.ConversationId, conversation.Type, messagesUpdate, conversationsUpdate, allUserConversationUpdates)
		}()
	}

	l.Logger.Infof("发送通知消息成功: conversationId=%s, messageId=%s, type=%d", in.ConversationId, messageId, in.MessageType)

	return &chat_rpc.SendNotificationMessageRes{
		MessageId: messageId,
	}, nil
}

// notifyNotificationUpdateGrouped 按会话分组推送通知更新（给该会话的所有用户推送）
func (l *SendNotificationMessageLogic) notifyNotificationUpdateGrouped(conversationId string, conversationType int, messagesUpdate, conversationsUpdate map[string]interface{}, allUserConversationUpdates []chat_utils.UserConversationUpdate) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("推送通知更新时发生panic: %v", r)
		}
	}()

	// 获取该会话的所有用户ID
	var recipientIds []string
	for _, update := range allUserConversationUpdates {
		if update.ConversationID == conversationId {
			recipientIds = append(recipientIds, update.UserID)
		}
	}

	// 构建按表分组的更新数据结构
	tableUpdates := []map[string]interface{}{
		messagesUpdate,      // Messages 更新
		conversationsUpdate, // Conversations 更新
	}

	// User_conversations 表更新
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
	fmt.Println("推送通知更新给会话成员: ", recipientIds)

	for _, recipientId := range recipientIds {
		// 一次性推送所有表的更新信息
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd.Hosts[0], wsCommandConst.CHAT_MESSAGE, messageType, "", recipientId, map[string]interface{}{
			"tableUpdates": tableUpdates, // 按表分组的更新数组
		}, conversationId)

		l.Logger.Infof("分组推送通知相关更新: recipient=%s, conversation=%s, tableUpdateCount=%d",
			recipientId, conversationId, len(tableUpdates))
	}
}
