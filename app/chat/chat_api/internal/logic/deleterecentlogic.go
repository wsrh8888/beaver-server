package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type DeleteRecentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewDeleteRecentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRecentLogic {
	return &DeleteRecentLogic{
		ctx:    ctx,
		logger: logger.New("delete_recent"),
		svcCtx: svcCtx,
	}
}

func (l *DeleteRecentLogic) DeleteRecent(req *types.DeleteRecentReq) (resp *types.DeleteRecentRes, err error) {
	// 假删除操作
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{"is_delete": true, "is_pinned": false}).Error

	if err == nil {
		l.logger.Info(model.LogMsg{
			Text: "删除会话成功",
			Data: map[string]interface{}{
				"userId":         req.UserID,
				"conversationId": req.ConversationID,
			},
		})
	}

	return nil, err
}
