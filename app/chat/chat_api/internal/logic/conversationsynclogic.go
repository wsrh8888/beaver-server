package logic

import (
	"context"
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 会话数据同步
func NewConversationSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationSyncLogic {
	return &ConversationSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConversationSyncLogic) ConversationSync(req *types.ConversationSyncReq) (resp *types.ConversationSyncRes, err error) {
	var conversations []chat_models.ChatConversationMeta

	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	// 查询数据 - 修复查询条件，ChatConversationMeta没有user_id字段
	err = l.svcCtx.DB.Where("version > ? AND version <= ?",
		req.FromVersion, req.ToVersion).
		Order("version ASC").
		Limit(limit + 1).
		Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话数据失败: %v", err)
		return nil, err
	}

	// 判断是否还有更多数据
	hasMore := len(conversations) > limit
	if hasMore {
		conversations = conversations[:limit]
	}

	// 转换为响应格式
	var conversationItems []types.ConversationSyncItem
	var nextVersion int64 = req.FromVersion

	for _, conv := range conversations {
		conversationItems = append(conversationItems, types.ConversationSyncItem{
			ConversationID: conv.ConversationID,
			Type:           conv.Type,
			LastReadSeq:    conv.LastReadSeq,
			Version:        conv.Version,
			CreateAt:       time.Time(conv.CreatedAt).Unix(),
			UpdateAt:       time.Time(conv.UpdatedAt).Unix(),
		})

		nextVersion = conv.Version
	}

	// 如果没有更多数据，nextVersion应该是toVersion+1
	if !hasMore {
		nextVersion = req.ToVersion + 1
	}

	resp = &types.ConversationSyncRes{
		Conversations: conversationItems,
		HasMore:       hasMore,
		NextVersion:   nextVersion,
	}

	l.Infof("会话数据同步完成，用户ID: %s, 返回会话数: %d, 还有更多: %v", req.UserID, len(conversationItems), hasMore)
	return resp, nil
}
