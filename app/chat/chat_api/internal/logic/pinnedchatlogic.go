package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinnedChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPinnedChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinnedChatLogic {
	return &PinnedChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PinnedChatLogic) PinnedChat(req *types.PinnedChatReq) (resp *types.PinnedChatRes, err error) {

	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversationModel{}).
		Where("user_id = ? AND conversation_id = ?", req, req.ConversationID).
		Updates(map[string]interface{}{"is_delete": false, "is_pinned": req.IsPinned}).Error

	return
}
