package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/chat/chat_utils"
	"beaver/common/models/ctype"

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

	// 更新会话级别的信息
	err = chat_utils.CreateOrUpdateConversation(l.svcCtx.DB, l.svcCtx.VersionGen, in.ConversationId, conversation.Type, nextSeq, in.Content)
	if err != nil {
		l.Logger.Errorf("更新会话信息失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 这里不返回错误，因为消息已经创建成功
	}

	// 获取会话中的所有用户，为每个用户更新会话关系
	var userConversations []chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).Find(&userConversations).Error
	if err != nil {
		l.Logger.Errorf("获取会话用户失败: conversationId=%s, error=%v", in.ConversationId, err)
		// 这里不返回错误，继续处理
	} else {
		// 为每个用户更新会话关系
		for _, userConv := range userConversations {
			err = chat_utils.UpdateUserConversation(l.svcCtx.DB, l.svcCtx.VersionGen, userConv.UserID, in.ConversationId, false)
			if err != nil {
				l.Logger.Errorf("更新用户会话关系失败: userId=%s, conversationId=%s, error=%v", userConv.UserID, in.ConversationId, err)
				// 这里不返回错误，继续处理其他用户
			}
		}
	}

	l.Logger.Infof("发送通知消息成功: conversationId=%s, messageId=%s, type=%d", in.ConversationId, messageId, in.MessageType)

	return &chat_rpc.SendNotificationMessageRes{
		MessageId: messageId,
	}, nil
}
