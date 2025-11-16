package logic

import (
	"context"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MuteChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置会话免打扰
func NewMuteChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteChatLogic {
	return &MuteChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MuteChatLogic) MuteChat(req *types.MuteChatReq) (resp *types.MuteChatRes, err error) {
	resp = &types.MuteChatRes{}

	// 获取下一个版本号
	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	// 更新会话免打扰状态和版本号
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_muted": req.IsMuted,
			"version":  version,
		}).Error
	if err != nil {
		l.Logger.Errorf("mute chat update failed: %v", err)
		return nil, err
	}

	return resp, nil
}
