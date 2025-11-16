package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type HideChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 隐藏/显示会话
func NewHideChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HideChatLogic {
	return &HideChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HideChatLogic) HideChat(req *types.HideChatReq) (resp *types.HideChatRes, err error) {
	resp = &types.HideChatRes{}

	// 获取下一个版本号
	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	// 更新会话隐藏状态和版本号
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_hidden": req.IsHidden,
			"version":   version,
		}).Error
	if err != nil {
		l.Logger.Errorf("hide chat update failed: %v", err)
		return nil, err
	}

	return resp, nil
}
