package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConversationVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationVersionsLogic {
	return &GetUserConversationVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConversationVersionsLogic) GetUserConversationVersions(in *chat_rpc.GetUserConversationVersionsReq) (*chat_rpc.GetUserConversationVersionsRes, error) {
	// 查询用户的所有会话设置
	var userConversations []chat_models.ChatUserConversation
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 如果提供了since参数，只返回更新时间大于since的记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话设置失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var userConversationVersions []*chat_rpc.UserConversationVersion
	for _, userConv := range userConversations {
		userConversationVersions = append(userConversationVersions, &chat_rpc.UserConversationVersion{
			ConversationId: userConv.ConversationID,
			Version:        userConv.Version,
		})
	}

	return &chat_rpc.GetUserConversationVersionsRes{
		UserConversationVersions: userConversationVersions,
	}, nil
}
