package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserConversationSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户会话关系数据同步
func NewUserConversationSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserConversationSyncLogic {
	return &UserConversationSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserConversationSyncLogic) UserConversationSync(req *types.UserConversationSyncReq) (resp *types.UserConversationSyncRes, err error) {
	var userConversations []chat_models.ChatUserConversation

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据
	err = l.svcCtx.DB.Where("user_id = ? AND version > ? AND version <= ?",
		req.UserID, req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话关系数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(userConversations) > limit
	if hasMore {
		userConversations = userConversations[:limit]
	}

	// 转换为响应格式
	var userConversationItems []types.UserConversationSyncItem
	var nextVersion int64 = req.FromVersion

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

		nextVersion = uc.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.UserConversationSyncRes{
		UserConversations: userConversationItems,
		HasMore:           hasMore,
		NextVersion:       nextVersion,
	}

	l.Infof("用户会话关系数据同步完成，用户ID: %s, 返回关系数: %d, 还有更多: %v", req.UserID, len(userConversationItems), hasMore)
	return resp, nil
}
