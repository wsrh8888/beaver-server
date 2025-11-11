package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DissolveConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDissolveConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DissolveConversationLogic {
	return &DissolveConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DissolveConversationLogic) DissolveConversation(in *chat_rpc.DissolveConversationReq) (*chat_rpc.DissolveConversationRes, error) {
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}

	// 检查会话是否存在
	var conversation chat_models.ChatConversationMeta
	if err := l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).First(&conversation).Error; err != nil {
		l.Errorf("会话不存在: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("会话不存在")
	}

	// 开启事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 确保在函数返回时处理事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除会话元数据
	if err := tx.Where("conversation_id = ?", in.ConversationId).Delete(&chat_models.ChatConversationMeta{}).Error; err != nil {
		l.Errorf("删除会话元数据失败: conversationId=%s, error=%v", in.ConversationId, err)
		tx.Rollback()
		return nil, errors.New("删除会话元数据失败")
	}

	// 2. 删除所有用户的会话关系记录
	if err := tx.Where("conversation_id = ?", in.ConversationId).Delete(&chat_models.ChatUserConversation{}).Error; err != nil {
		l.Errorf("删除用户会话关系失败: conversationId=%s, error=%v", in.ConversationId, err)
		tx.Rollback()
		return nil, errors.New("删除用户会话关系失败")
	}

	// 注意：消息记录通常保留作为历史记录，不在解散时删除
	// 如果需要删除消息，可以添加相应的逻辑

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	l.Infof("成功解散会话: conversationId=%s", in.ConversationId)

	return &chat_rpc.DissolveConversationRes{
		Success: true,
	}, nil
}
