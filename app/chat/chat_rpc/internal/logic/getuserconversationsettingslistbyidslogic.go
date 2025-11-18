package logic

import (
	"context"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationSettingsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConversationSettingsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationSettingsListByIdsLogic {
	return &GetUserConversationSettingsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConversationSettingsListByIdsLogic) GetUserConversationSettingsListByIds(in *chat_rpc.GetUserConversationSettingsListByIdsReq) (*chat_rpc.GetUserConversationSettingsListByIdsRes, error) {
	// 查询用户对指定会话的设置信息
	var userConversations []chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("user_id = ? AND conversation_id IN (?)", in.UserId, in.ConversationIds).Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话设置失败: %v", err)
		return nil, err
	}

	var userConversationSettings []*chat_rpc.UserConversationSettingListById
	for _, uc := range userConversations {
		userConversationSettings = append(userConversationSettings, &chat_rpc.UserConversationSettingListById{
			ConversationId: uc.ConversationID,
			Version:        uc.Version,
		})
	}

	return &chat_rpc.GetUserConversationSettingsListByIdsRes{
		UserConversationSettings: userConversationSettings,
	}, nil
}
