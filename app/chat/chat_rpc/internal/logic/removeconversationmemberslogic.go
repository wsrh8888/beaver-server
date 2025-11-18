package logic

import (
	"context"
	"errors"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveConversationMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveConversationMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveConversationMembersLogic {
	return &RemoveConversationMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveConversationMembersLogic) RemoveConversationMembers(in *chat_rpc.RemoveConversationMembersReq) (*chat_rpc.RemoveConversationMembersRes, error) {
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}

	if len(in.UserIds) == 0 {
		return nil, errors.New("用户列表不能为空")
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

	// 删除用户的会话关系记录（软删除或标记为隐藏）
	for _, userId := range in.UserIds {
		if err := tx.Model(&chat_models.ChatUserConversation{}).
			Where("conversation_id = ? AND user_id = ?", in.ConversationId, userId).
			Update("is_hidden", true).Error; err != nil {
			l.Errorf("移除用户会话关系失败: userId=%s, conversationId=%s, error=%v", userId, in.ConversationId, err)
			tx.Rollback()
			return nil, errors.New("移除用户会话关系失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	l.Infof("成功移除会话成员: conversationId=%s, userIds=%v", in.ConversationId, in.UserIds)

	return &chat_rpc.RemoveConversationMembersRes{
		Success: true,
	}, nil
}
