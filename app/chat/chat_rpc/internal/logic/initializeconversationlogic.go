package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitializeConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitializeConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitializeConversationLogic {
	return &InitializeConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitializeConversationLogic) InitializeConversation(in *chat_rpc.InitializeConversationReq) (*chat_rpc.InitializeConversationRes, error) {
	// 参数验证
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}

	if len(in.UserIds) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	if in.Type < 1 || in.Type > 4 {
		return nil, errors.New("无效的会话类型")
	}

	// 检查会话是否已存在
	var existingConversation chat_models.ChatConversationMeta
	err := l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).First(&existingConversation).Error
	if err == nil {
		// 会话已存在，直接返回
		l.Logger.Infof("会话已存在: conversationId=%s", in.ConversationId)
		return &chat_rpc.InitializeConversationRes{
			ConversationId: in.ConversationId,
		}, nil
	}

	// 创建会话元数据
	conversationMeta := chat_models.ChatConversationMeta{
		ConversationID: in.ConversationId,
		Type:           int(in.Type), // 转换类型
		MaxSeq:         0,            // 初始消息序列号
		Version:        1,            // 初始版本
	}
	// 注意：GORM会自动处理时间字段，无需手动赋值

	if err := l.svcCtx.DB.Create(&conversationMeta).Error; err != nil {
		l.Logger.Errorf("创建会话元数据失败: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("创建会话元数据失败")
	}

	// 为所有参与用户创建用户会话关系记录
	for _, userId := range in.UserIds {
		userConversation := chat_models.ChatUserConversation{
			UserID:         userId,
			ConversationID: in.ConversationId,
			IsPinned:       false,
			IsMuted:        false,
			UserReadSeq:    0,
			Version:        1,
		}
		// 注意：GORM会自动处理时间字段，无需手动赋值

		if err := l.svcCtx.DB.Create(&userConversation).Error; err != nil {
			l.Logger.Errorf("创建用户会话关系失败: userId=%s, conversationId=%s, error=%v", userId, in.ConversationId, err)
			return nil, errors.New("创建用户会话关系失败")
		}
	}

	l.Logger.Infof("初始化会话成功: conversationId=%s, type=%d, users=%v", in.ConversationId, in.Type, in.UserIds)

	return &chat_rpc.InitializeConversationRes{
		ConversationId: in.ConversationId,
	}, nil
}
