package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRecentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRecentLogic {
	return &DeleteRecentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRecentLogic) DeleteRecent(req *types.DeleteRecentReq) (resp *types.DeleteRecentRes, err error) {
	// 假删除操作
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{"is_delete": true, "is_pinned": false}).Error

	return nil, err
}
