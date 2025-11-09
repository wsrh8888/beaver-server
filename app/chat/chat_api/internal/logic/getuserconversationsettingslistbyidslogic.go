package logic

import (
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"context"
	"errors"
	"time"

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
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	if len(req.ConversationIds) == 0 {
		return &types.GetUserConversationSettingsListByIdsRes{
			UserConversationSettings: []types.UserConversationSyncItem{},
		}, nil
	}

	// 查询指定用户和会话ID的用户会话设置数据
	var userConversations []chat_models.ChatUserConversation
	err = l.svcCtx.DB.Where("user_id = ? AND conversation_id IN (?)", userId, req.ConversationIds).Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话设置失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var userConversationItems []types.UserConversationSyncItem
	for _, uc := range userConversations {
		userConversationItems = append(userConversationItems, types.UserConversationSyncItem{
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
		UserConversationSettings: userConversationItems,
	}, nil
}
