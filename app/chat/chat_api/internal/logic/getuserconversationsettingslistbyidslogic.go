package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserConversationSettingsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户会话设置数据
func NewGetUserConversationSettingsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationSettingsListByIdsLogic {
	return &GetUserConversationSettingsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserConversationSettingsListByIdsLogic) GetUserConversationSettingsListByIds(req *types.GetUserConversationSettingsListByIdsReq) (resp *types.GetUserConversationSettingsListByIdsRes, err error) {
	// 只查询用户会话设置表的数据
	var userConversations []chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("user_id = ? AND conversation_id IN (?)", req.UserID, req.ConversationIds).Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话设置失败: %v", err)
		return nil, err
	}

	// 转换数据库模型为API响应
	conversationSettings := make([]types.UserConversationSettingById, 0, len(userConversations))
	for _, uc := range userConversations {
		conversationSettings = append(conversationSettings, types.UserConversationSettingById{
			UserID:         uc.UserID,
			ConversationID: uc.ConversationID,
			IsHidden:       uc.IsHidden,
			IsPinned:       uc.IsPinned,
			IsMuted:        uc.IsMuted,
			UserReadSeq:    uc.UserReadSeq,
			Version:        uc.Version,
			CreateAt:       time.Time(uc.CreatedAt).Unix(),
			UpdateAt:       time.Time(uc.UpdatedAt).Unix(),
		})
	}

	return &types.GetUserConversationSettingsListByIdsRes{
		UserConversationSettings: conversationSettings,
	}, nil
}
